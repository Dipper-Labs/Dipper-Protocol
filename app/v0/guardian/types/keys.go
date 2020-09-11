package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/protocol"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	ModuleName   = protocol.GuardianModuleName
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

var (
	profilerKey = []byte{0x00}
)

func GetProfilerKey(addr sdk.AccAddress) []byte {
	return append(profilerKey, addr.Bytes()...)
}

func GetProfilersSubspaceKey() []byte {
	return profilerKey
}
