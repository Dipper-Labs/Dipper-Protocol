package types

import (
	"github.com/Dipper-Protocol/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBankBorrow{}, "dipperBank/BankBorrow", nil)
	cdc.RegisterConcrete(MsgBankRepay{}, "dipperBank/BankRepay", nil)
	cdc.RegisterConcrete(MsgBankDeposit{}, "dipperBank/BankDeposit", nil)
	cdc.RegisterConcrete(MsgBankWithdraw{}, "dipperBank/BankWithdraw", nil)
	cdc.RegisterConcrete(MsgSetOraclePrice{}, "dipperBank/SetOraclePrice", nil)
}
