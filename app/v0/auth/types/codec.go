package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/exported"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "dip/Account", nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "dip/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "dip/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "dip/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(StdTx{}, "dip/StdTx", nil)
}

// ModuleCdc - generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
