package keeper

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply/internal/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// DefaultCodespace from the supply module
var DefaultCodespace = types.ModuleName

// Keys for supply store
// Items are stored with the following key: values
// - 0x00: Supply
var (
	SupplyKey             = []byte{0x00}
	VestingStoreKeyPrefix = []byte{0x01}
)

// VestingStoreKey turn an address to key used to get it from the supply store
func VestingStoreKey(addr sdk.AccAddress) []byte {
	return append(VestingStoreKeyPrefix, addr.Bytes()...)
}
