package simapp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Dipper-Labs/Dipper-Protocol/app"
	v0 "github.com/Dipper-Labs/Dipper-Protocol/app/v0"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/simulation"
	"github.com/Dipper-Labs/Dipper-Protocol/baseapp"
	"github.com/Dipper-Labs/Dipper-Protocol/types/module"
)

func TestFullAppSimulation(t *testing.T) {
	config, db, dir, logger, skip, err := SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := app.NewDIPApp(logger, db, nil, true, FlagPeriodValue, baseapp.FauxMerkleMode())

	// run randomized simulation
	curProtocol := app.Engine.GetCurrentProtocol()
	cdc := curProtocol.GetCodec()
	smp := curProtocol.GetSimulationManager()
	sm, ok := smp.(*module.SimulationManager)
	require.True(t, ok)

	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(cdc, sm),
		SimulationOperations(app, cdc, config),
		v0.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		PrintStats(db)
	}
}
