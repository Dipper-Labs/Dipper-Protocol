package simapp

import (
	"flag"

	"github.com/Dipper-Labs/Dipper-Protocol/types/simulation"
)

func init() {
	GetSimulatorFlags()
}

// List of available flags for the simulator
var (
	FlagGenesisFileValue        string
	FlagParamsFileValue         string
	FlagExportParamsPathValue   string
	FlagExportParamsHeightValue int
	FlagExportStatePathValue    string
	FlagExportStatsPathValue    string
	FlagSeedValue               int64
	FlagInitialBlockHeightValue int
	FlagNumBlocksValue          int
	FlagBlockSizeValue          int
	FlagLeanValue               bool
	FlagCommitValue             bool
	FlagOnOperationValue        bool
	FlagAllInvariantsValue      bool

	FlagEnabledValue     bool
	FlagVerboseValue     bool
	FlagPeriodValue      uint
	FlagGenesisTimeValue int64
)

// GetSimulatorFlags gets the values of all the available simulation flags
func GetSimulatorFlags() {
	// config fields
	flag.StringVar(&FlagGenesisFileValue, "Genesis", "", "custom simulation genesis file; cannot be used with params file")
	flag.StringVar(&FlagParamsFileValue, "Params", "", "custom simulation params file which overrides any random params; cannot be used with genesis")
	flag.StringVar(&FlagExportParamsPathValue, "ExportParamsPath", "/tmp/ExportParamsPath", "custom file path to save the exported params JSON")
	flag.IntVar(&FlagExportParamsHeightValue, "ExportParamsHeight", 10, "height to which export the randomly generated params")
	flag.StringVar(&FlagExportStatePathValue, "ExportStatePath", "/tmp/ExportStatePath", "custom file path to save the exported app state JSON")
	flag.StringVar(&FlagExportStatsPathValue, "ExportStatsPath", "/tmp/ExportStatsPath", "custom file path to save the exported simulation statistics JSON")
	flag.Int64Var(&FlagSeedValue, "Seed", 42, "simulation random seed")
	flag.IntVar(&FlagInitialBlockHeightValue, "InitialBlockHeight", 1, "initial block to start the simulation")
	flag.IntVar(&FlagNumBlocksValue, "NumBlocks", 200, "number of new blocks to simulate from the initial block height")
	flag.IntVar(&FlagBlockSizeValue, "BlockSize", 200, "operations per block")
	flag.BoolVar(&FlagLeanValue, "Lean", false, "lean simulation log output")
	flag.BoolVar(&FlagCommitValue, "Commit", true, "have the simulation commit")
	flag.BoolVar(&FlagOnOperationValue, "SimulateEveryOperation", false, "run slow invariants every operation")
	flag.BoolVar(&FlagAllInvariantsValue, "PrintAllInvariants", false, "print all invariants if a broken invariant is found")

	// simulation flags
	flag.BoolVar(&FlagEnabledValue, "Enabled", true, "enable the simulation")
	flag.BoolVar(&FlagVerboseValue, "Verbose", false, "verbose log output")
	flag.UintVar(&FlagPeriodValue, "Period", 5, "run slow invariants only once every period assertions")
	flag.Int64Var(&FlagGenesisTimeValue, "GenesisTime", 0, "override genesis UNIX time instead of using a random UNIX time")
}

// NewConfigFromFlags creates a simulation from the retrieved values of the flags.
func NewConfigFromFlags() simulation.Config {
	return simulation.Config{
		GenesisFile:        FlagGenesisFileValue,
		ParamsFile:         FlagParamsFileValue,
		ExportParamsPath:   FlagExportParamsPathValue,
		ExportParamsHeight: FlagExportParamsHeightValue,
		ExportStatePath:    FlagExportStatePathValue,
		ExportStatsPath:    FlagExportStatsPathValue,
		Seed:               FlagSeedValue,
		InitialBlockHeight: FlagInitialBlockHeightValue,
		NumBlocks:          FlagNumBlocksValue,
		BlockSize:          FlagBlockSizeValue,
		Lean:               FlagLeanValue,
		Commit:             FlagCommitValue,
		OnOperation:        FlagOnOperationValue,
		AllInvariants:      FlagAllInvariantsValue,
	}
}
