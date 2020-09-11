package protocol

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// Protocol shows the expected behavior for any protocol version
type Protocol interface {
	GetVersion() uint64
	GetRouter() sdk.Router
	GetQueryRouter() sdk.QueryRouter
	GetAnteHandler() sdk.AnteHandler
	GetFeeRefundHandler() sdk.FeeRefundHandler
	GetInitChainer() sdk.InitChainer
	GetBeginBlocker() sdk.BeginBlocker
	GetEndBlocker() sdk.EndBlocker

	ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error)

	LoadContext()
	Init()
	GetCodec() *codec.Codec

	//for test
	SetRouter(sdk.Router)
	SetQueryRouter(sdk.QueryRouter)
	SetAnteHandler(anteHandler sdk.AnteHandler)
	SetInitChainer(sdk.InitChainer)

	//for simulation
	GetSimulationManager() interface{}
}
