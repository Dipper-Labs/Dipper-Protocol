package rest

import (
	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
	"github.com/gorilla/mux"
)

// REST query and parameter values
const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

// RegisterRoutes registers the auth module REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(
		"/auth/accounts/{address}", QueryAccountRequestHandlerFn(storeName, cliCtx),
	).Methods(MethodGet)

	r.HandleFunc(
		"/auth/params",
		queryParamsHandler(cliCtx),
	).Methods(MethodGet)

	r.HandleFunc(
		"/estimate_gas",
		EstimateGas(cliCtx),
	).Methods(MethodPost)
}

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/decode", DecodeTxRequestHandlerFn(cliCtx)).Methods("POST")
}
