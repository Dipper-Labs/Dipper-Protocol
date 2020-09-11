package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// RegisterCodec - register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgVerifyInvariant{}, "dip/MsgVerifyInvariant", nil)
}

// ModuleCdc - generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
