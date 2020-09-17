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
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// MaxProvisions: Maximum amount of mining, default value is 3.5e DIP
	if minter.CurrentProvisions.GTE(params.MaxProvisions) {
		return
	}

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

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	// update CurrentProvisions and save minter status to store
	minter.CurrentProvisions = minter.CurrentProvisions.Add(mintedCoin.Amount.ToDec())
	k.SetMinter(ctx, minter)

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
			sdk.NewAttribute(types.AttributeKeyCurrentProvisions, minter.CurrentProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
