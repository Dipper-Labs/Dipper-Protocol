package simapp

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	simapparams "github.com/Dipper-Labs/Dipper-Protocol/app/simapp/params"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	"github.com/Dipper-Labs/Dipper-Protocol/types/module"
	simtypes "github.com/Dipper-Labs/Dipper-Protocol/types/simulation"
)

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
func AppStateFn(cdc *codec.Codec, simManager *module.SimulationManager) simtypes.AppStateFn {

	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		genesisTimestamp = simtypes.RandTimestamp(r)
		chainID = config.ChainID

		appParams := make(simtypes.AppParams)
		appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)

		return appState, simAccs, chainID, genesisTimestamp
	}
}

// AppStateRandomizedFn creates calls each module's GenesisState generator function
// and creates the simulation params
func AppStateRandomizedFn(
	simManager *module.SimulationManager, r *rand.Rand, cdc *codec.Codec,
	accs []simtypes.Account, genesisTimestamp time.Time, appParams simtypes.AppParams,
) (json.RawMessage, []simtypes.Account) {
	numAccs := int64(len(accs))
	genesisState := NewDefaultGenesisState()

	// generate a random amount of initial stake coins and a random initial
	// number of bonded accounts
	var initialStake, numInitiallyBonded int64
	appParams.GetOrGenerate(
		cdc, simapparams.StakePerAccount, &initialStake, r,
		func(r *rand.Rand) {
			for initialStake < 1e17 {
				initialStake = r.Int63n(1e18)
			}
		},
	)
	appParams.GetOrGenerate(
		cdc, simapparams.InitiallyBondedValidators, &numInitiallyBonded, r,
		func(r *rand.Rand) { numInitiallyBonded = int64(r.Intn(300)) },
	)

	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}

	fmt.Printf(
		`Selected randomly generated parameters for simulated genesis:
{
  stake_per_account: "%d",
  initially_bonded_validators: "%d"
}
`, initialStake, numInitiallyBonded,
	)

	simState := &module.SimulationState{
		AppParams:    appParams,
		Cdc:          cdc,
		Rand:         r,
		GenState:     genesisState,
		Accounts:     accs,
		InitialStake: initialStake,
		NumBonded:    numInitiallyBonded,
		GenTimestamp: genesisTimestamp,
	}

	simManager.GenerateGenesisStates(simState)

	appState, err := cdc.MarshalJSON(genesisState)
	if err != nil {
		panic(err)
	}

	return appState, accs
}
