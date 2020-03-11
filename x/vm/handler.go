package vm

import (
	"fmt"
	"github.com/Dipper-Protocol/x/vm/keeper"
	"github.com/Dipper-Protocol/x/vm/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/Dipper-Protocol/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgContract:
			return handleMsgContract(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
			//return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgContract(ctx sdk.Context, msg MsgContract, k Keeper) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	gasLimit := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	_, res, err2 := DoStateTransition(ctx, msg, k, gasLimit, false)
	//TODO temporary, need to fix
	if err2 != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Data: res.Data, GasUsed: res.GasUsed, Events: ctx.EventManager().Events()}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	k.StateDB.UpdateAccounts() //update account balance for fee refund when create/call contract
	k.StateDB.WithContext(ctx).Commit(true)
	return []abci.ValidatorUpdate{}
}
