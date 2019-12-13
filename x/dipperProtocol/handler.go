package dipperProtocol

import (
	"fmt"
	"github.com/Dipper-Protocol/x/dipperProtocol/internal/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "dipperProtocol" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
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
			errMsg := fmt.Sprintf("Unrecognized dipperProtocol Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	// Checks if the the bid price is greater than the price paid by the current owner
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg MsgDeleteName) sdk.Result {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return types.ErrNameDoesNotExist(types.DefaultCodespace).Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}

	keeper.DeleteWhois(ctx, msg.Name)
	return sdk.Result{}
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