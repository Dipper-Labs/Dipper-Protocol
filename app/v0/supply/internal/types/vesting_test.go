package types

import (
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CalculateVestingAmount(t *testing.T) {
	vestingAmount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000)))

	//[startTime, endTime]
	vestings := []Vesting{
		{sdk.AccAddress{0x1}, vestingAmount, 100, 1100},
	}

	vesting := CalculateVestingAmount(90, vestings)
	require.Equal(t, vesting, vestingAmount)

	vesting = CalculateVestingAmount(100, vestings)
	require.Equal(t, vesting, vestingAmount)

	vesting = CalculateVestingAmount(1100, vestings)
	require.Equal(t, vesting, sdk.NewCoins())

	vesting = CalculateVestingAmount(1200, vestings)
	require.Equal(t, vesting, sdk.NewCoins())

	vesting = CalculateVestingAmount(101, vestings)
	require.Equal(t, vesting.AmountOf(sdk.DefaultBondDenom).Int64(), sdk.NewInt(999).Int64())

	vesting = CalculateVestingAmount(700, vestings)
	require.Equal(t, vesting.AmountOf(sdk.DefaultBondDenom).Int64(), sdk.NewInt(400).Int64())

	vesting = CalculateVestingAmount(1099, vestings)
	require.Equal(t, vesting.AmountOf(sdk.DefaultBondDenom).Int64(), sdk.NewInt(1).Int64())

	//[0, endTime]
	vestings = []Vesting{
		{sdk.AccAddress{0x1}, vestingAmount, 0, 1000},
	}

	vesting = CalculateVestingAmount(1, vestings)
	require.Equal(t, vesting, vestingAmount)

	vesting = CalculateVestingAmount(500, vestings)
	require.Equal(t, vesting, vestingAmount)

	vesting = CalculateVestingAmount(999, vestings)
	require.Equal(t, vesting, vestingAmount)

	vesting = CalculateVestingAmount(1000, vestings)
	require.Equal(t, vesting, sdk.NewCoins())

	vesting = CalculateVestingAmount(1001, vestings)
	require.Equal(t, vesting, sdk.NewCoins())
}
