package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddProfiler{}, "dip/guardian/MsgAddProfiler", nil)
	cdc.RegisterConcrete(MsgDeleteProfiler{}, "dip/guardian/MsgDeleteProfiler", nil)
	cdc.RegisterConcrete(Guardian{}, "dip/guardian/Guardian", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
