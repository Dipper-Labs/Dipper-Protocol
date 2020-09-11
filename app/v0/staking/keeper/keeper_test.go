package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking/types"
)

func TestParams(t *testing.T) {
	ctx, _, keeper, _ := CreateTestInput(t, false, 0)
	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	resParams := keeper.GetParams(ctx)
	resParams.NextExtendingTime = expParams.NextExtendingTime
	require.True(t, expParams.Equal(resParams))

	//modify a params, save, and retrieve
	expParams.MaxValidators = 777
	keeper.SetParams(ctx, expParams)
	resParams = keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}
