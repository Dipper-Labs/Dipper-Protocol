package dipperBank

import (
	"fmt"
	"github.com/Dipper-Protocol/x/dipperBank/internal/types"
	"strconv"

	sdk "github.com/Dipper-Protocol/types"
)

// NewHandler returns a handler for "dipperBank" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgBankBorrow:
			return handleMsgBankBorrow(ctx, keeper, msg)
		case MsgBankRepay:
			return handleMsgBankRepay(ctx, keeper, msg)
		case MsgBankDeposit:
			return handleMsgBankDeposit(ctx, keeper, msg)
		case MsgBankWithdraw:
			return handleMsgBankWithdraw(ctx, keeper, msg)
		case MsgSetOraclePrice:
			return handleMsgSetOraclePrice(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized dipperBank Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgBankBorrow(ctx sdk.Context, keeper Keeper, msg MsgBankBorrow) sdk.Result{
	err := keeper.BankBorrow(ctx, msg.Amount, msg.Symbol, msg.Owner)
	if err != nil {
		return types.ErrNotEnoughTokenForBorrow(types.DefaultCodespace).Result()
	}
	keeper.CoinKeeper.SendCoins(ctx, DipperBankAddress, msg.Owner, msg.Amount)
	return sdk.Result{}
}

func handleMsgBankRepay(ctx sdk.Context, keeper Keeper, msg MsgBankRepay) sdk.Result{
	err := keeper.BankRepay(ctx, msg.Amount, msg.Symbol, msg.Owner)
	if err != nil {
		return types.ErrTooMuchAmountToRepay(types.DefaultCodespace).Result()
	}
	keeper.CoinKeeper.SendCoins(ctx, msg.Owner, DipperBankAddress, msg.Amount)
	return sdk.Result{}
}

func handleMsgBankDeposit(ctx sdk.Context, keeper Keeper, msg MsgBankDeposit) sdk.Result{
	keeper.BankDeposit(ctx, msg.Amount, msg.Symbol, msg.Owner)
	keeper.CoinKeeper.SendCoins(ctx, msg.Owner, DipperBankAddress, msg.Amount)
	return sdk.Result{}
}

func handleMsgBankWithdraw(ctx sdk.Context, keeper Keeper, msg MsgBankWithdraw) sdk.Result{
	err := keeper.BankWithdraw(ctx, msg.Amount, msg.Symbol, msg.Owner)
	if err != nil {
		return types.ErrNotEnoughAmountCoinForWithdraw(types.DefaultCodespace).Result()
	}
	keeper.CoinKeeper.SendCoins(ctx, DipperBankAddress, msg.Owner, msg.Amount)
	return sdk.Result{}
}

func handleMsgSetOraclePrice(ctx sdk.Context, keeper Keeper, msg MsgSetOraclePrice) sdk.Result{
	//TODO add authority who can set oracle price.
	//if msg.Owner.Equals(){
	//
	//}
	price, err := strconv.ParseInt(msg.Price, 10, 64)
	if err != nil {
		return sdk.ErrUnknownRequest("invalid price").Result()
	}
	keeper.SetOraclePrice(ctx, msg.Symbol, price)
	return sdk.Result{}
}