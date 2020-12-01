package vm

import (
	"fmt"
	"math/big"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/vm/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	DefaultBlockGasLimit = sdk.Gas(100000000)
)

// StateTransition defines data to transitionDB in vm
type StateTransition struct {
	Sender    sdk.AccAddress
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
	StateDB   *types.CommitStateDB
}

func (st StateTransition) CanTransfer(acc sdk.AccAddress, amount *big.Int) bool {
	return st.StateDB.GetBalance(acc).Cmp(amount) >= 0
}

func (st StateTransition) Transfer(from, to sdk.AccAddress, amount *big.Int) {
	st.StateDB.SubBalance(from, amount)
	st.StateDB.AddBalance(to, amount)
}

func (st StateTransition) GetHashFn(header abci.Header) func() sdk.Hash {
	return func() sdk.Hash {
		var res = sdk.Hash{}
		blockID := header.GetLastBlockId()
		res.SetBytes(blockID.GetHash())
		return res
	}
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context, k Keeper) (*big.Int, *sdk.Result, error) {
	logger := k.Logger(ctx)

	evmCtx := Context{
		CanTransfer: st.CanTransfer,
		Transfer:    st.Transfer,
		GetHash:     st.GetHashFn(ctx.BlockHeader()),
		Origin:      st.Sender,
		CoinBase:    ctx.BlockHeader().ProposerAddress,
		Time:        sdk.NewInt(ctx.BlockHeader().Time.Unix()).BigInt(),
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	gasLimitForVM := DefaultBlockGasLimit
	if ctx.BlockGasMeter() != nil {
		gasLimitForVM = ctx.BlockGasMeter().Limit()
	}

	if !ctx.Simulate {
		gasLimitForVM = ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed() - GasReservedForOutOfVm
		if gasLimitForVM >= ctx.GasMeter().Limit() {
			return nil, &sdk.Result{Data: nil, GasUsed: ctx.GasMeter().GasConsumed()}, sdk.ErrOutOfGas("")
		}
	}
	evmCtx.GasLimit = gasLimitForVM

	curGasMeter := ctx.GasMeter()
	gasMeterForEvm := sdk.NewInfiniteGasMeter()

	vmParams := k.GetParams(ctx) // will consume gas
	st.StateDB.UpdateAccounts()  // will consume gas

	cfg := Config{
		OpConstGasConfig:          &vmParams.VMOpGasParams,
		ContractCreationGasConfig: &vmParams.VMContractCreationGasParams,
		MaxCodeSize:               vmParams.MaxCodeSize,
		MaxCallCreateDepth:        vmParams.MaxCallCreateDepth,
	}
	evm := NewEVM(evmCtx, st.StateDB.WithContext(ctx.WithGasMeter(gasMeterForEvm)), cfg)

	var (
		ret         []byte
		leftOverGas uint64
		vmerr       error
	)

	if st.Recipient.Empty() {
		ret, _, leftOverGas, vmerr = evm.Create(st.Sender, st.Payload, gasLimitForVM, st.Amount.BigInt())
		logger.Info(fmt.Sprintf("create contract, consumed gas = %v, leftOverGas = %v, vm err = %v ", gasLimitForVM-leftOverGas, leftOverGas, vmerr))
	} else {
		ret, leftOverGas, vmerr = evm.Call(st.Sender, st.Recipient, st.Payload, gasLimitForVM, st.Amount.BigInt())
		if vmerr == ErrExecutionReverted {
			reason := "null"
			if len(ret) > 4 {
				reason = string(ret[4:])
			}
			logger.Info(fmt.Sprintf("VM revert error, reason provided by the contract: %s", reason))
		}

		logger.Info(fmt.Sprintf("call contract, ret = %x, consumed gas = %v, leftOverGas = %v, vm err = %v", ret, gasLimitForVM-leftOverGas, leftOverGas, vmerr))
	}

	vmGasUsed := gasLimitForVM - leftOverGas

	if vmerr != nil {
		ctx.EventManager().Clear()
		ctx.WithGasMeter(curGasMeter).GasMeter().ConsumeGas(vmGasUsed, "VM execution consumption")
		return nil, &sdk.Result{Data: ret, GasUsed: curGasMeter.GasConsumed() + vmGasUsed}, vmerr
	}

	st.StateDB.Finalise(true)

	// comsume vm gas
	ctx.WithGasMeter(curGasMeter).GasMeter().ConsumeGas(vmGasUsed, "VM execution consumption")

	return nil, &sdk.Result{Data: ret, GasUsed: ctx.GasMeter().GasConsumed()}, nil
}

func DoStateTransition(ctx sdk.Context, msg types.MsgContract, k Keeper, readonly bool) (*big.Int, *sdk.Result, error) {
	st := StateTransition{
		Sender:    msg.From,
		Recipient: msg.To,
		Payload:   msg.Payload,
		Amount:    msg.Amount.Amount,
		StateDB:   k.StateDB.WithContext(ctx).WithTxHash(tmhash.Sum(ctx.TxBytes())),
	}

	if readonly {
		ctx.Simulate = true
	}

	if !ctx.Simulate && ctx.GasMeter().Limit() == 0 {
		return nil, &sdk.Result{Data: nil}, ErrWrongCtx
	}

	if ctx.Simulate {
		st.StateDB = types.NewStateDB(k.StateDB).WithContext(ctx)
	}

	return st.TransitionCSDB(ctx, k)
}
