package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// ModuleCdc - generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ParameterChangeProposal{}, "dip/ParameterChangeProposal", nil)
}
