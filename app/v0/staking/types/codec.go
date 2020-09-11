package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// RegisterCodec - register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "dip/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "dip/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "dip/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "dip/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "dip/MsgBeginRedelegate", nil)
}

// ModuleCdc - generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
