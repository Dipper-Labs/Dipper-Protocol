package rest

import (
	"github.com/gorilla/mux"

	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
