package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Dipper-Labs/Dipper-Protocol/app/protocol"
	v0 "github.com/Dipper-Labs/Dipper-Protocol/app/v0"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	v1 "github.com/Dipper-Labs/Dipper-Protocol/app/v1"
	"github.com/Dipper-Labs/Dipper-Protocol/baseapp"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/Dipper-Labs/Dipper-Protocol/types/module"
	"github.com/Dipper-Labs/Dipper-Protocol/version"
)

const (
	appName = "dip"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.dipcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.dipd")
)

// DIPApp extends BaseApp
type DIPApp struct {
	*baseapp.BaseApp
}

// Codec returns the current protocol codec
func (app *DIPApp) Codec() *codec.Codec {
	return app.Engine.GetCurrentProtocol().GetCodec()
}

// BeginBlocker abci
func (app *DIPApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.BeginBlock(req)
}

// EndBlocker abci
func (app *DIPApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.EndBlock(req)
}

// InitChainer - custom logic for initialization
func (app *DIPApp) InitChainer(_ sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	return app.InitChain(req)
}

// TODO: check
// ModuleAccountAddrs returns all the module account addresses
func (app *DIPApp) ModuleAccountAddrs() map[string]bool {
	return nil
}

// SimulationManager implements the SimulationApp interface
func (app *DIPApp) SimulationManager() *module.SimulationManager {
	smp := app.Engine.GetCurrentProtocol().GetSimulationManager()
	sm, ok := smp.(*module.SimulationManager)
	if !ok {
		return nil
	}

	return sm
}

// NewDIPApp returns a reference to an initialized DIPApp
func NewDIPApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp)) *DIPApp {
	baseApp := baseapp.NewBaseApp(appName, logger, db, baseAppOptions...)

	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	mainStoreKey := protocol.Keys[protocol.MainStoreKey]
	protocolKeeper := sdk.NewProtocolKeeper(mainStoreKey)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	baseApp.SetProtocolEngine(&engine)
	baseApp.MountKVStores(protocol.Keys)
	baseApp.MountTransientStores(protocol.TKeys)

	var app = &DIPApp{baseApp}

	// set hook function postEndBlocker
	baseApp.PostEndBlocker = app.postEndBlocker

	if loadLatest {
		err := app.LoadLatestVersion(mainStoreKey)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, app.DeliverTx, invCheckPeriod, nil))
	engine.Add(v1.NewProtocolV1(1, logger, protocolKeeper, app.DeliverTx, invCheckPeriod, nil))

	loaded, current := engine.LoadCurrentProtocol(app.GetCms().GetKVStore(mainStoreKey))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!, to upgrade dipd", current))
	}
	logger.Info(fmt.Sprintf("launch app with protocol version: %d", current))

	// set txDeocder
	app.SetTxDecoder(auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec()))

	return app
}

func MakeLatestCodec() *codec.Codec {
	return v1.MakeCodec()
}

func (app *DIPApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.Keys[protocol.MainStoreKey])
}

// hook function for BaseApp's EndBlock(upgrade)
func (app *DIPApp) postEndBlocker(res *abci.ResponseEndBlock) {
	appVersion := app.Engine.GetCurrentVersion()
	for _, event := range res.Events {
		if event.Type == sdk.AppVersionEvent {
			for _, attr := range event.Attributes {
				if string(attr.Key) == sdk.AppVersionEvent {
					appVersion, _ = strconv.ParseUint(string(attr.Value), 10, 64)
					break
				}
			}

			break
		}
	}

	if appVersion <= app.Engine.GetCurrentVersion() {
		return
	}

	success := app.Engine.Activate(appVersion)
	if success {
		app.SetTxDecoder(auth.DefaultTxDecoder(app.Engine.GetCurrentProtocol().GetCodec()))
		return
	}

	app.Log(fmt.Sprintf("activate version from %d to %d failed, please upgrade your app", app.Engine.GetCurrentVersion(), appVersion))
	os.Exit(0)
}

// ExportAppStateAndValidators exports the state of application for a genesis file
func (app *DIPApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, jailWhiteList)
}
