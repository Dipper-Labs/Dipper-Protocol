package types

const (
	// ModuleName is the name of the module
	ModuleName = "dipperBank"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// Each Module has its own name
	DipperBank = "dipperBank"


)

//Dipper-Bank
//this is the token pool
type DUser = string
type DSymbol = string
type DBill = int64
type DPrice = int64