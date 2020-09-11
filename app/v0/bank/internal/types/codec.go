package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// RegisterCodec - Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "dip/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "dip/MsgMultiSend", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
