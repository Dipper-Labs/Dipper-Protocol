package mint

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/mint/internal/types"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

func deleteFinishedVestings(ctx sdk.Context, k Keeper, vestings []supply.Vesting) {
	for _, vesting := range vestings {
		if ctx.BlockTime().Unix() > vesting.EndTime {
			k.RemoveVesting(ctx, vesting.Address)
		}
	}
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// calculate supplyExcludingVesting
	totalStakingSupply := k.StakingTokenSupply(ctx)
	vestings := k.GetAllVestings(ctx)
	vestingAmount := supply.CalculateVestingAmount(ctx.BlockTime().Unix(), vestings)
	supplyExcludingVesting := totalStakingSupply.Sub(vestingAmount.AmountOf(sdk.DefaultBondDenom))
	deleteFinishedVestings(ctx, k, vestings)

	// recalculate inflation rate
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, supplyExcludingVesting)
	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
