package vm

import (
	"encoding/hex"
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/types"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	sdkerrors "github.com/Dipper-Labs/Dipper-Protocol/types/errors"
)

const (
	DefaultVMGasLimit = 100000000
)

func NewQuerier(k keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryState:
			return queryState(ctx, req, k)
		case types.QueryCode:
			return queryCode(ctx, path, k)
		case types.QueryStorage:
			return queryStorage(ctx, path, k)
		case types.QueryTxLogs:
			return queryTxLogs(ctx, path, k)
		case types.EstimateGas, types.QueryCall:
			return simulateStateTransition(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParameters(ctx sdk.Context, k keeper.Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryState(ctx sdk.Context, req abci.RequestQuery, k keeper.Keeper) (res []byte, err error) {
	var params types.QueryStateParams
	err = codec.Cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, err
	}

	stateObjects := k.StateDB.WithContext(ctx).ExportStateObjects(params)
	res, err = json.Marshal(stateObjects)
	return
}

func queryCode(ctx sdk.Context, path []string, k keeper.Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	code := k.GetCode(ctx, addr)

	return code, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper keeper.Keeper) ([]byte, error) {
	addr, _ := sdk.AccAddressFromBech32(path[1])
	key := sdk.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	bRes := types.QueryStorageResult{Value: val}
	res, err := codec.MarshalJSONIndent(keeper.Cdc, bRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryTxLogs(ctx sdk.Context, path []string, keeper keeper.Keeper) ([]byte, error) {
	txHash := sdk.HexToHash(path[1])
	logs := keeper.GetLogs(ctx, txHash)

	bRes := types.QueryLogsResult{Logs: logs}
	res, err := codec.MarshalJSONIndent(keeper.Cdc, bRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func simulateStateTransition(ctx sdk.Context, req abci.RequestQuery, k keeper.Keeper) ([]byte, error) {
	var msg types.MsgContract
	codec.Cdc.UnmarshalJSON(req.Data, &msg)

	_, result, err := DoStateTransition(ctx, msg, k, true)

	if err == nil {
		bRes := types.SimulationResult{Gas: result.GasUsed, Res: hex.EncodeToString(result.Data)}
		res, err := codec.MarshalJSONIndent(k.Cdc, bRes)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return res, nil
	}

	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "StateTransition faileds")
}
