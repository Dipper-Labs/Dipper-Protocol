package v1

import (
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

func NewDefaultGenesisState() sdk.GenesisState {
	return ModuleBasics.DefaultGenesis()
}
