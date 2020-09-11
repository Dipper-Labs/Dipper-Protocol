package baseapp

import (
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// nolint - Mostly for testing
func (app *BaseApp) Check(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	return app.runTx(runTxModeCheck, nil, tx)
}

// nolint - full tx execution
func (app *BaseApp) Simulate(txBytes []byte, tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	return app.runTx(runTxModeSimulate, txBytes, tx)
}

// nolint
func (app *BaseApp) Deliver(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	return app.runTx(runTxModeDeliver, nil, tx)
}

// Context with current {check, deliver}State of the app
// used by tests
func (app *BaseApp) NewContext(isCheckTx bool, header abci.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, header, true, app.logger).
			WithMinGasPrices(app.minGasPrices)
	}

	return sdk.NewContext(app.deliverState.ms, header, false, app.logger)
}

func (app *BaseApp) SetTxDecoder(txDecoder sdk.TxDecoder) {
	app.txDecoder = txDecoder
}
