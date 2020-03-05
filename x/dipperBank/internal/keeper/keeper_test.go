package keeper

import (
	"fmt"
	"testing"
)

func TestKeeper_SetOraclePrice(t *testing.T) {
	ctx, _, _, _, _, keeper := CreateTestInput(t, false, 1000)
	billbank := keeper.GetBillBank(ctx)
	oracle, _ := billbank.GetOracle()
	fmt.Println("1", oracle)
	keeper.SetOraclePrice(ctx, "eth", 150000000)
	keeper.SetOraclePrice(ctx, "dai", 150000000)
	keeper.SetOraclePrice(ctx, "btc", 150000000)
	billbank = keeper.GetBillBank(ctx)
	oracle, _ = billbank.GetOracle()
	fmt.Println("3", oracle)

}