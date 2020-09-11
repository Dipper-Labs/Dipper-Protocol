package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// ModuleCdc - generic sealed codec to be used throughout this module
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)

	cdc.RegisterConcrete(MsgSubmitProposal{}, "dip/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "dip/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "dip/MsgVote", nil)

	cdc.RegisterConcrete(TextProposal{}, "dip/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, "dip/SoftwareUpgradeProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)
}
