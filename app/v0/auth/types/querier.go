package types

import (
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount = "account"
	QueryParams  = "params"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address sdk.AccAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}
