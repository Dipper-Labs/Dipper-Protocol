package types

const (
	// ModuleName is the name of the module
	ModuleName = "dipperProtocol"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// Each Module has its own name
	DipperBank = "dipperBank"


)

//Dipper-Bank
//this is the token pool
type tUser = string
type tSymbol = string
type tBill = int64
type tPrice = int64