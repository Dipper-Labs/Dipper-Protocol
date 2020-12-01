package vm

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/common"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/types"
)

const (
	ModuleName               = types.ModuleName
	StoreKey                 = types.StoreKey
	RouterKey                = types.RouterKey
	QuerierRoute             = types.QuerierRoute
	DefaultParamspace        = keeper.DefaultParamspace
	EventTypeContractCreated = types.EventTypeContractCreated
	AttributeKeyAddress      = types.AttributeKeyAddress
)

type (
	Keeper        = keeper.Keeper
	AccountKeeper = types.AccountKeeper
	MsgContract   = types.MsgContract
	CommitStateDB = types.CommitStateDB
	Log           = types.Log
	Params        = types.Params

	GenesisState = types.GenesisState
)

var (
	// functions aliases
	NewKeeper        = keeper.NewKeeper
	NewCommitStateDB = types.NewCommitStateDB
	NewMsgContract   = types.NewMsgContract
	NewParams        = types.NewParams
	DefaultParams    = types.DefaultParams

	CreateAddress  = common.CreateAddress
	CreateAddress2 = common.CreateAddress2
	GenPayload     = cli.GenPayload
	CodeFromFile   = cli.CodeFromFile

	ValidateGenesis = types.ValidateGenesis

	ErrOutOfGas                 = types.ErrOutOfGas
	ErrCodeStoreOutOfGas        = types.ErrCodeStoreOutOfGas
	ErrDepth                    = types.ErrDepth
	ErrTraceLimitReached        = types.ErrTraceLimitReached
	ErrInsufficientBalance      = types.ErrInsufficientBalance
	ErrContractAddressCollision = types.ErrContractAddressCollision
	ErrNoCompatibleInterpreter  = types.ErrNoCompatibleInterpreter
	ErrEmptyInputs              = types.ErrEmptyInputs
	ErrNoCodeExist              = types.ErrNoCodeExist
	ErrMaxCodeSizeExceeded      = types.ErrMaxCodeSizeExceeded
	ErrWriteProtection          = types.ErrWriteProtection
	ErrReturnDataOutOfBounds    = types.ErrReturnDataOutOfBounds
	ErrExecutionReverted        = types.ErrExecutionReverted
	ErrInvalidJump              = types.ErrInvalidJump
	ErrGasUintOverflow          = types.ErrGasUintOverflow
	ErrNoPayload                = types.ErrNoPayload
	ErrWrongCtx                 = types.ErrWrongCtx

	// variable aliases
	ModuleCdc = types.ModuleCdc
)
