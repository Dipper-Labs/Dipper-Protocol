package bank

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/bank/internal/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/bank/internal/types"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/vm"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	sdkerrors "github.com/Dipper-Labs/Dipper-Protocol/types/errors"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)

		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgSend) (*sdk.Result, error) {
	if !k.GetSendEnabled(ctx) {
		return nil, ErrSendDisabled
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", msg.ToAddress)
	}

	if k.IsContractAccount(ctx, msg.ToAddress) {
		var contractMsg vm.MsgContract
		contractMsg.From = msg.FromAddress
		contractMsg.To = msg.ToAddress
		contractMsg.Amount = sdk.Coin{Denom: sdk.NativeTokenName, Amount: msg.Amount.AmountOf(sdk.NativeTokenName)}
		vm.HandleMsgContract(ctx, contractMsg, k.GetVMKeeper())
	}

	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgMultiSend) (*sdk.Result, error) {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return nil, ErrSendDisabled
	}

	for _, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", out.Address)
		}
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
