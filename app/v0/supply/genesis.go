package supply

import (
	autypes "github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply/internal/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

func CalculateVestingAmount(blockTime int64, vestings []Vesting) sdk.Coins {
	amt := sdk.NewCoins()
	for _, vesting := range vestings {
		if blockTime >= vesting.EndTime {
			continue
		}

		if vesting.StartTime > 0 {
			if blockTime > vesting.StartTime {
				x := vesting.EndTime - blockTime
				y := vesting.EndTime - vesting.StartTime
				s := sdk.NewDec(x).Quo(sdk.NewDec(y))

				vestingCoins := sdk.NewCoins()
				for _, ovc := range vesting.Amount {
					vestingAmt := ovc.Amount.ToDec().Mul(s).RoundInt()
					vestingCoins = append(vestingCoins, sdk.NewCoin(ovc.Denom, vestingAmt))
				}

				amt = amt.Add(vestingCoins)
			}
		} else if vesting.StartTime == 0 {
			amt = amt.Add(vesting.Amount)
		}
	}
	return amt
}

func vestingInfoFromAccount(acc autypes.Account) (isVestingAccount bool, vesting types.Vesting) {
	switch accObj := acc.(type) {
	case *autypes.DelayedVestingAccount:
		vesting.Address = accObj.Address
		vesting.Amount = accObj.OriginalVesting
		vesting.StartTime = 0
		vesting.EndTime = accObj.EndTime
		return true, vesting

	case *autypes.ContinuousVestingAccount:
		vesting.Address = accObj.Address
		vesting.Amount = accObj.OriginalVesting
		vesting.StartTime = accObj.StartTime
		vesting.EndTime = accObj.EndTime
		return true, vesting
	}

	return false, vesting
}

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper Keeper, ak types.AccountKeeper, data GenesisState) {
	// manually set the total supply based on accounts if not provided
	if data.Supply.Empty() {
		var totalSupply sdk.Coins
		ak.IterateAccounts(ctx,
			func(acc autypes.Account) (stop bool) {
				totalSupply = totalSupply.Add(acc.GetCoins())

				isVestingAccount, vesting := vestingInfoFromAccount(acc)
				if isVestingAccount {
					keeper.SetVesting(ctx, vesting)
				}

				return false
			},
		)

		data.Supply = totalSupply
	}

	keeper.SetSupply(ctx, types.NewSupply(data.Supply))
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetSupply(ctx).GetTotal())
}

// ValidateGenesis performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return types.NewSupply(data.Supply).ValidateBasic()
}
