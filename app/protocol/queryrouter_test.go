package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var testQuerier = func(_ sdk.Context, _ []string, _ abci.RequestQuery) (res []byte, err error) {
	return nil, nil
}

func TestQueryRouter(t *testing.T) {
	qr := NewQueryRouter()

	// require panic on invalid route
	require.Panics(t, func() {
		qr.AddRoute("*", testQuerier)
	})

	qr.AddRoute("testRoute", testQuerier)
	q := qr.Route("testRoute")
	require.NotNil(t, q)

	// require panic on duplicate route
	require.Panics(t, func() {
		qr.AddRoute("testRoute", testQuerier)
	})
}
