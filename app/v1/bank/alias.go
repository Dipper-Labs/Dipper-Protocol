// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
package bank

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/bank/internal/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/bank/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	RegisterCodec          = types.RegisterCodec
	ErrNoInputs            = types.ErrNoInputs
	ErrNoOutputs           = types.ErrNoOutputs
	ErrInputOutputMismatch = types.ErrInputOutputMismatch
	ErrSendDisabled        = types.ErrSendDisabled
	NewBaseKeeper          = keeper.NewBaseKeeper
	NewInput               = types.NewInput
	NewOutput              = types.NewOutput
	ParamKeyTable          = types.ParamKeyTable
	NewMsgSend             = types.NewMsgSend

	// variable aliases
	ModuleCdc                = types.ModuleCdc
	ParamStoreKeySendEnabled = types.ParamStoreKeySendEnabled
)

type (
	BaseKeeper   = keeper.BaseKeeper // ibc module depends on this
	Keeper       = keeper.Keeper
	MsgSend      = types.MsgSend
	MsgMultiSend = types.MsgMultiSend
	Input        = types.Input
	Output       = types.Output
)
