package types

// MinterKey - the one key to use for the keeper store
var MinterKey = []byte{0x00}

// nolint
const (
	// ModuleName defines the name of the module
	ModuleName = "mint"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters        = "parameters"
	QueryInflation         = "inflation"
	QueryAnnualProvisions  = "annual_provisions"
	QueryCurrentProvisions = "current_provisions"
)
