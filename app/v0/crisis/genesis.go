package crisis

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/crisis/internal/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/crisis/internal/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// new crisis genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	keeper.SetConstantFee(ctx, data.ConstantFee)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	constantFee := keeper.GetConstantFee(ctx)
	return types.NewGenesisState(constantFee)
}
