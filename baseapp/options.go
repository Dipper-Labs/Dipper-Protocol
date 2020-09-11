// nolint: golint
package baseapp

import (
	"fmt"
	"io"

	dbm "github.com/tendermint/tm-db"

	"github.com/Dipper-Labs/Dipper-Protocol/store"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// File for storing in-package BaseApp optional functions,
// for options that need access to non-exported fields of the BaseApp

// SetPruning sets a pruning option on the multistore associated with the app
func SetPruning(opts sdk.PruningOptions) func(*BaseApp) {
	return func(bap *BaseApp) { bap.cms.SetPruning(opts) }
}

// SetMinGasPrices returns an option that sets the minimum gas prices on the app.
func SetMinGasPrices(gasPricesStr string) func(*BaseApp) {
	gasPrices, err := sdk.ParseDecCoins(gasPricesStr)
	if err != nil {
		panic(fmt.Sprintf("invalid minimum gas prices: %v", err))
	}

	return func(bap *BaseApp) { bap.setMinGasPrices(gasPrices) }
}

// SetHaltHeight returns a BaseApp option function that sets the halt height.
func SetHaltHeight(height uint64) func(*BaseApp) {
	return func(bap *BaseApp) { bap.setHaltHeight(height) }
}

func (app *BaseApp) SetName(name string) {
	if app.sealed {
		panic("SetName() on sealed BaseApp")
	}
	app.name = name
}

// SetAppVersion sets the application's version string.
func (app *BaseApp) SetAppVersion(v string) {
	if app.sealed {
		panic("SetAppVersion() on sealed BaseApp")
	}
	app.appVersion = v
}

func (app *BaseApp) SetDB(db dbm.DB) {
	if app.sealed {
		panic("SetDB() on sealed BaseApp")
	}
	app.db = db
}

func (app *BaseApp) SetCMS(cms store.CommitMultiStore) {
	if app.sealed {
		panic("SetEndBlocker() on sealed BaseApp")
	}
	app.cms = cms
}

func (app *BaseApp) SetAddrPeerFilter(pf sdk.PeerFilter) {
	if app.sealed {
		panic("SetAddrPeerFilter() on sealed BaseApp")
	}
	app.addrPeerFilter = pf
}

func (app *BaseApp) SetIDPeerFilter(pf sdk.PeerFilter) {
	if app.sealed {
		panic("SetIDPeerFilter() on sealed BaseApp")
	}
	app.idPeerFilter = pf
}

func (app *BaseApp) SetFauxMerkleMode() {
	if app.sealed {
		panic("SetFauxMerkleMode() on sealed BaseApp")
	}
	app.fauxMerkleMode = true
}

func FauxMerkleMode() func(*BaseApp) {
	return func(app *BaseApp) { app.SetFauxMerkleMode() }
}

// SetCommitMultiStoreTracer sets the store tracer on the BaseApp's underlying
// CommitMultiStore.
func (app *BaseApp) SetCommitMultiStoreTracer(w io.Writer) {
	app.cms.SetTracer(w)
}
