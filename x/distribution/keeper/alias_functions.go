package keeper

import (
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/x/distribution/types"
	"github.com/Dipper-Protocol/x/supply/exported"
)

// get outstanding rewards
func (k Keeper) GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val sdk.ValAddress) sdk.DecCoins {
	return k.GetValidatorOutstandingRewards(ctx, val)
}

// get the community coins
func (k Keeper) GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins {
	return k.GetFeePool(ctx).CommunityPool
}

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
