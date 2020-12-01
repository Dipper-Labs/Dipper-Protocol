package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v1/vm/types"
	"github.com/Dipper-Labs/Dipper-Protocol/types/rest"
	"github.com/gorilla/mux"

	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/vm/storage/{addr}/{key}",
		getStorageFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/vm/%s", types.EstimateGas),
		estimateGasFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/vm/%s/{addr}", types.QueryCode),
		getCodeFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/vm/logs/{txId}",
		getLogFn(cliCtx),
	).Methods("GET")

	// Get the current staking parameter values
	r.HandleFunc(
		"/vm/parameters",
		paramsHandlerFn(cliCtx),
	).Methods("GET")
}

func queryStorage(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["addr"]
		key := vars["key"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/vm/%s/%s/%s", types.QueryStorage, addr, key)
		res, height, err := cliCtx.Query(route)
		if err != nil {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func estimateGas(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params types.MsgContractQuery
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &params) {
			return
		}

		if params.From == nil || params.Payload == nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "bad request")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		d, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			return
		}

		route := fmt.Sprintf("custom/vm/%s", types.EstimateGas)
		res, height, err := cliCtx.QueryWithData(route, d)
		if err != nil {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getCode(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["addr"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/vm/%s/%s", types.QueryCode, addr)
		res, height, err := cliCtx.Query(route)
		if err != nil {
			return
		}

		dst := make([]byte, 2*len(res))
		hex.Encode(dst, res)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, string(dst))
	}
}

func getLog(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		txID := vars["txId"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/vm/logs/%s", txID)
		res, height, err := cliCtx.Query(route)
		if err != nil {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getParams(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.Query("custom/vm/params")
		if err != nil {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getStorageFn(cliCtx context.CLIContext) http.HandlerFunc {
	return queryStorage(cliCtx)
}

func estimateGasFn(cliCtx context.CLIContext) http.HandlerFunc {
	return estimateGas(cliCtx)
}

func getCodeFn(cliCtx context.CLIContext) http.HandlerFunc {
	return getCode(cliCtx)
}

func getLogFn(cliCtx context.CLIContext) http.HandlerFunc {
	return getLog(cliCtx)
}

// HTTP request handler to query the staking params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return getParams(cliCtx)
}
