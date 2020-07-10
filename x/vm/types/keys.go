package types

const (
	// ModuleName is the name of the vm module
	ModuleName = "vm"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	CodeKey       = StoreKey + "_code"
	LogKey        = StoreKey + "_log"
	StoreDebugKey = StoreKey + "_debug"


	// QuerierRoute is the querier route for the vm module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the vm module
	RouterKey = ModuleName
)

var (
	LogIndexKey = []byte("logIndexKey")
)
