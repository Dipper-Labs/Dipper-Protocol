package types

import (
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/x/supply/exported"
)

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}
