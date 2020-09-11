package simulation

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/genaccounts/internal/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/Dipper-Labs/Dipper-Protocol/types/module"
)

func RandomGenesisAccounts(simState *module.SimulationState) (genesisAccs types.GenesisAccounts) {
	for _, acc := range simState.Accounts {
		bacc := auth.NewBaseAccountWithAddress(acc.Address)
		coins := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, simState.InitialStake))
		bacc.SetCoins(coins)
		gacc := types.NewGenesisAccount(&bacc)
		genesisAccs = append(genesisAccs, gacc)
	}

	return genesisAccs
}
