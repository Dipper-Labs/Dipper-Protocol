package distribution

// DONTCOVER

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/client/rest"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/simulation"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/types"
	stakingkeeper "github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking/keeper"
	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/Dipper-Labs/Dipper-Protocol/types/module"
	simtypes "github.com/Dipper-Labs/Dipper-Protocol/types/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the distr module.
type AppModuleBasic struct{}

// Name returns the distr module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, StoreKey)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd returns the root query command for the distr module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

// AppModule implements an application module for the distr module.
type AppModule struct {
	AppModuleBasic
	keeper          Keeper
	supplyKeeper    types.SupplyKeeper
	akForSimulation keeper.AccountKeeper // for simulation
	skForSimulation stakingkeeper.Keeper // for simulation
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		supplyKeeper:   supplyKeeper,
	}
}

func (am *AppModule) WithAccountKeeper(ak keeper.AccountKeeper) *AppModule {
	am.akForSimulation = ak
	return am
}

func (am *AppModule) WithStakingKeeper(sk stakingkeeper.Keeper) *AppModule {
	am.skForSimulation = sk
	return am
}

// Name returns the distr module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the distr module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the distr module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the distr module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns no sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the distr module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the distr
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the distr module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the distr module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState creates a randomized GenState of the distr module.
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// WeightedOperations returns the all the distr module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.keeper, am.akForSimulation, am.skForSimulation)
}

// ProposalContents returns all the distr content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return simulation.ProposalContents(am.keeper)
}

// RandomizedParams creates randomized distr param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return simulation.ParamChanges(r)
}
