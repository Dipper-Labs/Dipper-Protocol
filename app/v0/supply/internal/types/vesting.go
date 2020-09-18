package types

import sdk "github.com/Dipper-Labs/Dipper-Protocol/types"

type Vesting struct {
	Address   sdk.AccAddress `json:"address"`
	Amount    sdk.Coins      `json:"Amount"`
	StartTime int64          `json:"start_time"`
	EndTime   int64          `json:"end_time"`
}

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
			} else {
				amt = amt.Add(vesting.Amount)
			}
		} else if vesting.StartTime == 0 {
			amt = amt.Add(vesting.Amount)
		}
	}

	return amt
}
