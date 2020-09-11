package simapp

import (
	v0 "github.com/Dipper-Labs/Dipper-Protocol/app/v0"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() sdk.GenesisState {
	return v0.ModuleBasics.DefaultGenesis()
}
