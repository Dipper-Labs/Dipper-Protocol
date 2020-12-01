package types

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/exported"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a modules/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAllAccounts(ctx sdk.Context) []exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)

	IterateAccounts(ctx sdk.Context, process func(exported.Account) bool)
}
