package slashing

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/exported"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
}
