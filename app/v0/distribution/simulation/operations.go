package simulation

import (
	"fmt"
	"math/rand"

	"github.com/Dipper-Labs/Dipper-Protocol/app/simapp/helpers"
	simappparams "github.com/Dipper-Labs/Dipper-Protocol/app/simapp/params"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/types"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/simulation"
	stakingkeeper "github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/baseapp"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	simtypes "github.com/Dipper-Labs/Dipper-Protocol/types/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgSetWithdrawAddress          = "op_weight_msg_set_withdraw_address"
	OpWeightMsgWithdrawDelegationReward    = "op_weight_msg_withdraw_delegation_reward"
	OpWeightMsgWithdrawValidatorCommission = "op_weight_msg_withdraw_validator_commission"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simtypes.AppParams, cdc *codec.Codec, k keeper.Keeper, ak keeper.AccountKeeper, sk stakingkeeper.Keeper) simulation.WeightedOperations {

	var weightMsgSetWithdrawAddress int
	appParams.GetOrGenerate(cdc, OpWeightMsgSetWithdrawAddress, &weightMsgSetWithdrawAddress, nil,
		func(_ *rand.Rand) {
			weightMsgSetWithdrawAddress = simappparams.DefaultWeightMsgSetWithdrawAddress
		},
	)

	var weightMsgWithdrawDelegationReward int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawDelegationReward, &weightMsgWithdrawDelegationReward, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawDelegationReward = simappparams.DefaultWeightMsgWithdrawDelegationReward
		},
	)

	var weightMsgWithdrawValidatorCommission int
	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawValidatorCommission, &weightMsgWithdrawValidatorCommission, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawValidatorCommission = simappparams.DefaultWeightMsgWithdrawValidatorCommission
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgSetWithdrawAddress,
			SimulateMsgSetWithdrawAddress(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgWithdrawDelegationReward,
			SimulateMsgWithdrawDelegatorReward(ak, k, sk),
		),
		//simulation.NewWeightedOperation(
		//	weightMsgWithdrawValidatorCommission,
		//	SimulateMsgWithdrawValidatorCommission(ak, k, sk),
		//),
	}
}

// SimulateMsgSetWithdrawAddress generates a MsgSetWithdrawAddress with random values.
func SimulateMsgSetWithdrawAddress(ak keeper.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetWithdrawAddress, "app invalid"), nil, nil
		}

		if !k.GetWithdrawAddrEnabled(ctx) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetWithdrawAddress, "withdrawal is not enabled"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		simToAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		coins := account.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, coins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetWithdrawAddress, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgSetWithdrawAddress(simAccount.Address, simToAccount.Address)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawDelegatorReward generates a MsgWithdrawDelegatorReward with random values.
func SimulateMsgWithdrawDelegatorReward(ak keeper.AccountKeeper, k keeper.Keeper, sk stakingkeeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawDelegatorReward, "app invalid"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		delegations := sk.GetAllDelegatorDelegations(ctx, simAccount.Address)
		if len(delegations) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawDelegatorReward, "number of delegators equal 0"), nil, nil
		}

		delegation := delegations[r.Intn(len(delegations))]

		validator := sk.Validator(ctx, delegation.GetValidatorAddr())
		if validator == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawDelegatorReward, "validator is nil"), nil, fmt.Errorf("validator %s not found", delegation.GetValidatorAddr())
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := account.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawDelegatorReward, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgWithdrawDelegatorReward(simAccount.Address, validator.GetOperator())

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawValidatorCommission generates a MsgWithdrawValidatorCommission with random values.
func SimulateMsgWithdrawValidatorCommission(ak keeper.AccountKeeper, k keeper.Keeper, sk stakingkeeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app interface{}, ctx sdk.Context, accs []simtypes.Account, chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		a := baseapp.DereferenceBaseApp(app)
		if a == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawValidatorCommission, "app invalid"), nil, nil
		}

		validator, ok := stakingkeeper.RandomValidator(r, sk, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawValidatorCommission, "random validator is not ok"), nil, nil
		}

		simAccount, found := simtypes.FindAccount(accs, sdk.AccAddress(validator.GetOperator()))
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawValidatorCommission, "could not find account"), nil, fmt.Errorf("validator %s not found", validator.GetOperator())
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := account.GetCoins()

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawValidatorCommission, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgWithdrawValidatorCommission(validator.GetOperator())

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = a.Deliver(tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
