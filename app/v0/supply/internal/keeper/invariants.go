package keeper

import (
	"fmt"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/exported"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply/internal/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// RegisterInvariants register all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "total-supply", TotalSupply(k))
}

// AllInvariants runs all invariants of the supply module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return TotalSupply(k)(ctx)
	}
}

// TotalSupply checks that the total supply reflects all the coins held in accounts
func TotalSupply(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var expectedTotal sdk.Coins
		supply := k.GetSupply(ctx)

		k.ak.IterateAccounts(ctx, func(acc exported.Account) bool {
			expectedTotal = expectedTotal.Add(acc.GetCoins())
			return false
		})

		broken := !expectedTotal.IsEqual(supply.GetTotal())

		return sdk.FormatInvariant(types.ModuleName, "total supply",
			fmt.Sprintf(
				"\tsum of accounts coins: %v\n"+
					"\tsupply.Total:          %v\n",
				expectedTotal, supply.GetTotal())), broken
	}
}
