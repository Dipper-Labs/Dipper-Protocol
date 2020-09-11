package types

import (
	authexported "github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/exported"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}
