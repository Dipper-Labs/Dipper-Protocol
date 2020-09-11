package slashing

// DONTCOVER

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/slashing/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/slashing/client/rest"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/slashing/types"
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

// app module basics object
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
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
	rest.RegisterRoutes(ctx, rtr)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the slashing module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	keeper          Keeper
	stakingKeeper   types.StakingKeeper
	akForSimulation AccountKeeper        // for simulation
	skForSimulation stakingkeeper.Keeper // for simulation
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, stakingKeeper types.StakingKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		stakingKeeper:  stakingKeeper,
	}
}

func (am *AppModule) WithAccountKeeper(ak AccountKeeper) *AppModule {
	am.akForSimulation = ak
	return am
}

func (am *AppModule) WithStakingKeeper(sk stakingkeeper.Keeper) *AppModule {
	am.skForSimulation = sk
	return am
}

// module name
func (AppModule) Name() string {
	return ModuleName
}

// register invariants
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// module message route name
func (AppModule) Route() string {
	return RouterKey
}

// module handler
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.stakingKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// for simulation
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	RandomizedGenState(simState)
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return WeightedOperations(simState.AppParams, simState.Cdc, am.akForSimulation, am.keeper, am.skForSimulation)
}

func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return ParamChanges(r)
}
