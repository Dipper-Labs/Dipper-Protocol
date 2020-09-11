package baseapp

import (
	"fmt"
	"os"
	"sort"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	sdkerrors "github.com/Dipper-Labs/Dipper-Protocol/types/errors"
	"github.com/Dipper-Labs/Dipper-Protocol/version"
)

func (app *BaseApp) InitChain(req abci.RequestInitChain) (res abci.ResponseInitChain) {
	// Stash the consensus params in the cms main store and memoize.
	if req.ConsensusParams != nil {
		app.setConsensusParams(req.ConsensusParams)
		app.storeConsensusParams(req.ConsensusParams)
	}

	initHeader := abci.Header{ChainID: req.ChainId, Time: req.Time}

	// initialize the deliver state and check state with a correct header
	app.setDeliverState(initHeader)
	app.setCheckState(initHeader)

	initChainer := app.Engine.GetCurrentProtocol().GetInitChainer()
	if initChainer == nil {
		return
	}

	// add block gas meter for any genesis transactions (allow infinite gas)
	app.deliverState.ctx = app.deliverState.ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

	res = initChainer(app.deliverState.ctx, req)

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(res.Validators) {
			panic(
				fmt.Errorf(
					"len(RequestInitChain.Validators) != len(GenesisValidators) (%d != %d)",
					len(req.Validators), len(res.Validators),
				),
			)
		}

		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(res.Validators))

		for i, val := range res.Validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("genesisValidators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}
	// There may be some application state in the genesis file, so always init the metrics.
	//app.Engine.GetCurrentProtocol().InitMetrics(app.cms)

	return res
}

func (app *BaseApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	lastCommitID := app.cms.LastCommitID()

	return abci.ResponseInfo{
		AppVersion:       version.AppVersion,
		Data:             app.name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

func (app *BaseApp) SetOption(req abci.RequestSetOption) (res abci.ResponseSetOption) {
	return
}

func (app *BaseApp) FilterPeerByAddrPort(info string) abci.ResponseQuery {
	if app.addrPeerFilter != nil {
		return app.addrPeerFilter(info)
	}
	return abci.ResponseQuery{}
}

func (app *BaseApp) FilterPeerByID(info string) abci.ResponseQuery {
	if app.idPeerFilter != nil {
		return app.idPeerFilter(info)
	}
	return abci.ResponseQuery{}
}

func splitPath(requestPath string) (path []string) {
	path = strings.Split(requestPath, "/")
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}
	return path
}

func (app *BaseApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
	if app.cms.TracingEnabled() {
		app.cms.SetTracingContext(sdk.TraceContext(
			map[string]interface{}{"blockHeight": req.Header.Height},
		))
	}

	if err := app.validateHeight(req); err != nil {
		panic(err)
	}

	// Initialize the DeliverTx state. If this is the first block, it should
	// already be initialized in InitChain. Otherwise app.deliverState will be
	// nil, since it is reset on Commit.
	if app.deliverState == nil {
		app.setDeliverState(req.Header)
	} else {
		// In the first block, app.deliverState.ctx will already be initialized
		// by InitChain. Context is now updated with Header information.
		app.deliverState.ctx = app.deliverState.ctx.
			WithBlockHeader(req.Header).
			WithBlockHeight(req.Header.Height)
	}

	var gasMeter sdk.GasMeter
	if maxGas := app.getMaximumBlockGas(); maxGas > 0 {
		gasMeter = sdk.NewGasMeter(maxGas)
	} else {
		gasMeter = sdk.NewInfiniteGasMeter()
	}

	app.deliverState.ctx = app.deliverState.ctx.WithBlockGasMeter(gasMeter)

	beginBlocker := app.Engine.GetCurrentProtocol().GetBeginBlocker()
	if beginBlocker != nil {
		res = beginBlocker(app.deliverState.ctx, req)
	}

	// set the signed validators for addition to context in deliverTx
	app.voteInfos = req.LastCommitInfo.GetVotes()

	return res
}

func (app *BaseApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	if app.deliverState.ms.TracingEnabled() {
		app.deliverState.ms = app.deliverState.ms.SetTracingContext(nil).(sdk.CacheMultiStore)
	}

	endBlocker := app.Engine.GetCurrentProtocol().GetEndBlocker()
	if endBlocker != nil {
		res = endBlocker(app.deliverState.ctx, req)
	}

	if app.PostEndBlocker != nil {
		app.PostEndBlocker(&res)
	}

	return res
}

// CheckTx implements the ABCI interface. It runs the "basic checks" to see
// whether or not a transaction can possibly be executed, first decoding and then
// the ante handler (which checks signatures/fees/ValidateBasic).
//
// NOTE:CheckTx does not run the actual Msg handler function(s).
func (app *BaseApp) CheckTx(req abci.RequestCheckTx) (res abci.ResponseCheckTx) {
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return sdkerrors.ResponseCheckTx(err, 0, 0)
	}

	var mode runTxMode

	switch {
	case req.Type == abci.CheckTxType_New:
		mode = runTxModeCheck

	case req.Type == abci.CheckTxType_Recheck:
		mode = runTxModeReCheck

	default:
		panic(fmt.Sprintf("unknown RequestCheckTx type: %s", req.Type))
	}

	gInfo, result, err := app.runTx(mode, req.Tx, tx)
	if err != nil {
		return sdkerrors.ResponseCheckTx(err, gInfo.GasWanted, gInfo.GasUsed)
	}

	return abci.ResponseCheckTx{
		GasWanted: int64(gInfo.GasWanted),
		GasUsed:   int64(gInfo.GasUsed),
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}

func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, nil, 0, 0, err.Error())
	}

	gInfo, result, err := app.runTx(runTxModeDeliver, req.Tx, tx)

	if err != nil {
		log := err.Error()
		if result != nil {
			log = result.Log
		}

		var data []byte
		if result != nil {
			data = append(data, result.Data...)
		}

		return sdkerrors.ResponseDeliverTx(err, data, gInfo.GasWanted, gInfo.GasUsed, log)
	}

	return abci.ResponseDeliverTx{
		GasWanted: int64(gInfo.GasWanted),
		GasUsed:   int64(gInfo.GasUsed),
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}

// Commit implements the ABCI interface. It will commit all state that exists in
// the deliver state's multi-store and includes the resulting commit ID in the
// returned abci.ResponseCommit. Commit will set the check state based on the
// latest header and reset the deliver state. Also, if a non-zero halt height is
// defined in config, Commit will execute a deferred function call to check
// against that height and gracefully halt if it matches the latest committed
// height.
func (app *BaseApp) Commit() (res abci.ResponseCommit) {
	header := app.deliverState.ctx.BlockHeader()

	// write the Deliver state and commit the MultiStore
	app.deliverState.ms.Write()
	commitID := app.cms.Commit()
	app.logger.Debug("Commit synced", "commit", fmt.Sprintf("%X", commitID))

	// Reset the Check state to the latest committed.
	//
	// NOTE: This is safe because Tendermint holds a lock on the mempool for
	// Commit. Use the header from this latest block.
	app.setCheckState(header)

	// empty/reset the deliver state
	app.deliverState = nil

	defer func() {
		if app.haltHeight > 0 && uint64(header.Height) == app.haltHeight {
			app.Log("halting node per configuration height:%d", app.haltHeight)
			os.Exit(0)
		}
	}()

	return abci.ResponseCommit{
		Data: commitID.Hash,
	}
}

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *BaseApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	path := splitPath(req.Path)
	if len(path) == 0 {
		sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "no query path provided"))
	}

	switch path[0] {
	case "app":
		return handleQueryApp(app, path, req)

	case "store":
		return handleQueryStore(app, path, req)

	case "p2p":
		return handleQueryP2P(app, path, req)

	case "custom":
		return handleQueryCustom(app, path, req)
	}

	return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query path"))
}

func handleQueryApp(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(path) >= 2 {
		switch path[1] {
		case "simulate":
			txBytes := req.Data

			tx, err := app.txDecoder(txBytes)
			if err != nil {
				return sdkerrors.QueryResult(sdkerrors.Wrap(err, "failed to decode tx"))
			}

			gInfo, _, _ := app.Simulate(txBytes, tx)

			if req.Height == 0 {
				req.Height = app.LastBlockHeight()
			}

			return abci.ResponseQuery{
				Codespace: sdkerrors.RootCodespace,
				Height:    req.Height,
				Value:     codec.Cdc.MustMarshalBinaryLengthPrefixed(gInfo.GasUsed),
			}
		case "version":
			return abci.ResponseQuery{
				Codespace: sdkerrors.RootCodespace,
				Height:    req.Height,
				Value:     []byte(app.appVersion),
			}

		default:
			return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query: %s", path))
		}
	}

	return sdkerrors.QueryResult(
		sdkerrors.Wrap(
			sdkerrors.ErrUnknownRequest,
			"expected second parameter to be either 'simulate' or 'version', neither was present",
		),
	)
}

func handleQueryStore(app *BaseApp, path []string, req abci.RequestQuery) abci.ResponseQuery {
	queryable, ok := app.cms.(sdk.Queryable)
	if !ok {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "multistore doesn't support queries"))
	}

	req.Path = "/" + strings.Join(path[1:], "/")

	if req.Height == 0 {
		req.Height = app.LastBlockHeight()
	}

	if req.Height <= 1 && req.Prove {
		return sdkerrors.QueryResult(
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			),
		)
	}

	resp := queryable.Query(req)
	resp.Height = req.Height

	return resp
}

func handleQueryP2P(app *BaseApp, path []string, _ abci.RequestQuery) (res abci.ResponseQuery) {
	// "/p2p" prefix for p2p queries
	if len(path) >= 4 {
		cmd, typ, arg := path[1], path[2], path[3]
		switch cmd {
		case "filter":
			switch typ {
			case "addr":
				return app.FilterPeerByAddrPort(arg)

			case "id":
				return app.FilterPeerByID(arg)
			}

		default:
			return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "expected second parameter to be 'filter'"))
		}
	}

	return sdkerrors.QueryResult(
		sdkerrors.Wrap(
			sdkerrors.ErrUnknownRequest, "expected path is p2p filter <addr|id> <parameter>",
		),
	)
}

func handleQueryCustom(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	// path[0] should be "custom" because "/custom" prefix is required for keeper
	// queries.
	//
	// The queryRouter routes using path[1]. For example, in the path
	// "custom/gov/proposal", queryRouter routes using "gov".
	if len(path) < 2 || path[1] == "" {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "no route for custom query specified"))
	}

	router := app.Engine.GetCurrentProtocol().GetQueryRouter()
	querier := router.Route(path[1])
	if querier == nil {
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "no custom querier found for route %s", path[1]))
	}

	if req.Height == 0 {
		req.Height = app.LastBlockHeight()
	}

	if req.Height <= 1 && req.Prove {
		return sdkerrors.QueryResult(
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			),
		)
	}

	cacheMS, err := app.cms.CacheMultiStoreWithVersion(req.Height)
	if err != nil {
		return sdkerrors.QueryResult(
			sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"failed to load state at height %d; %s (latest height: %d)", req.Height, err, app.LastBlockHeight(),
			),
		)
	}

	// cache wrap the commit-multistore for safety
	ctx := sdk.NewContext(
		cacheMS, app.checkState.ctx.BlockHeader(), true, app.logger,
	).WithMinGasPrices(app.minGasPrices)

	// Passes the rest of the path as an argument to the querier.
	//
	// For example, in the path "custom/gov/proposal/test", the gov querier gets
	// []string{"proposal", "test"} as the path.
	resBytes, queryErr := querier(ctx, path[2:], req)
	if queryErr != nil {
		space, code, log := sdkerrors.ABCIInfo(err, false)
		return abci.ResponseQuery{
			Code:      code,
			Codespace: space,
			Log:       log,
			Height:    req.Height,
		}
	}

	return abci.ResponseQuery{
		Height: req.Height,
		Value:  resBytes,
	}
}

// Log format
func (app *BaseApp) Log(format string, a ...interface{}) {
	app.logger.Info(fmt.Sprintf(format, a...))
}
