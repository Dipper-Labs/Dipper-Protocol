package types

import (
	"github.com/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Protocol/types"
	authtypes "github.com/Dipper-Protocol/x/auth/types"
	stakingtypes "github.com/Dipper-Protocol/x/staking/types"
)

// ModuleCdc defines a generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// TODO: abstract genesis transactions registration back to staking
// required for genesis transactions
func init() {
	ModuleCdc = codec.New()
	stakingtypes.RegisterCodec(ModuleCdc)
	authtypes.RegisterCodec(ModuleCdc)
	sdk.RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
