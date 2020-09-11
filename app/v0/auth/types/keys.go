package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/protocol"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	// ModuleName defines the name of the module
	ModuleName = protocol.AuthModuleName

	// StoreKey is string representation of the store key for auth
	StoreKey = ModuleName

	// FeeCollectorName the root string for the fee collector account address
	FeeCollectorName = "fee_collector"

	// QuerierRoute is the querier route for acc
	QuerierRoute = ModuleName

	RefundKey = "refund_fee"
)

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	// param key for global account number
	GlobalAccountNumberKey = []byte("globalAccountNumber")
)

// AddressStoreKey turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}
