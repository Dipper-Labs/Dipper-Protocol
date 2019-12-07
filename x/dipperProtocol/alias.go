package dipperProtocol

import (
	"github.com/Dipper-Protocol/x/dipperProtocol/internal/keeper"
	"github.com/Dipper-Protocol/x/dipperProtocol/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper        = keeper.NewKeeper
	NewQuerier       = keeper.NewQuerier
	NewMsgBuyName    = types.NewMsgBuyName
	NewMsgSetName    = types.NewMsgSetName
	NewMsgDeleteName = types.NewMsgDeleteName
	NewWhois         = types.NewWhois
	ModuleCdc        = types.ModuleCdc
	RegisterCodec    = types.RegisterCodec


	NewBankBorrow  = types.NewMsgBankBorrow
	NewBankRepay = types.NewMsgBankRepay
	NewBankDeposit = types.NewMsgBankDeposit
	NewBankWithdraw = types.NewMsgBankWithdraw
)

type (
	Keeper          = keeper.Keeper
	MsgSetName      = types.MsgSetName
	MsgBuyName      = types.MsgBuyName
	MsgDeleteName   = types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	Whois           = types.Whois


	MsgBankBorrow = types.MsgBankBorrow
	MsgBankRepay = types.MsgBankRepay
	MsgBankDeposit = types.MsgBankDeposit
	MsgBankWithdraw = types.MsgBankWithdraw
	BillBank = types.BillBank
	TokenPool = types.TokenPool
)
