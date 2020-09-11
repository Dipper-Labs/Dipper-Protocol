package upgrade

// DONTCOVER

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/upgrade/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/upgrade/client/rest"
	upgtypes "github.com/Dipper-Labs/Dipper-Protocol/app/v0/upgrade/types"
	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/Dipper-Labs/Dipper-Protocol/types/module"
)

// check the implementation of the interface
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is a struct of app module basics object
type AppModuleBasic struct{}

// Name returns module name
func (a AppModuleBasic) Name() string {
	return upgtypes.ModuleName
}

// RegisterCodec registers module codec
func (a AppModuleBasic) RegisterCodec(*codec.Codec) {
}

// DefaultGenesis returns default genesis state
func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	d, _ := json.Marshal(DefaultGenesisState())
	return d
}

// ValidateGenesis validates genesis
func (a AppModuleBasic) ValidateGenesis(d json.RawMessage) error {
	var gs GenesisState
	return json.Unmarshal(d, &gs)
}

// RegisterRESTRoutes register rest routes
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the transaction commands for this module
func (a AppModuleBasic) GetTxCmd(*codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd gets the root query command of the upgrade module
func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(upgtypes.ModuleName, cdc)
}

// AppModule is a struct of app module
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object for upgrade module
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}

// InitGenesis initializes module genesis
func (a AppModule) InitGenesis(ctx sdk.Context, d json.RawMessage) []types.ValidatorUpdate {
	var gs GenesisState
	err := json.Unmarshal(d, &gs)
	if err != nil {
		panic(err)
	}
	InitGenesis(ctx, a.keeper, gs)
	return nil
}

// ExportGenesis exports module genesis
func (a AppModule) ExportGenesis(sdk.Context) json.RawMessage {
	d, err := json.Marshal(ExportGenesis())
	if err != nil {
		panic(err)
	}
	return d
}

// RegisterInvariants performs a no-op.
func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	// no op
}

// Route returns module message route name
func (a AppModule) Route() string {
	return upgtypes.RouterKey
}

// NewHandler returns an sdk.Handler for the upgrade module.
func (a AppModule) NewHandler() sdk.Handler {
	return nil
}

// QuerierRoute returns module querier route name
func (a AppModule) QuerierRoute() string {
	return upgtypes.QuerierRoute
}

// NewQuerierHandler returns the auth module sdk.Querier.
func (a AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

// BeginBlock returns the begin blocker for the upgrade module.
func (a AppModule) BeginBlock(sdk.Context, types.RequestBeginBlock) {
	panic("implement me")
}

// EndBlock returns the end blocker for the upgrade module. It returns no validator
// updates.
func (a AppModule) EndBlock(ctx sdk.Context, b types.RequestEndBlock) []types.ValidatorUpdate {
	EndBlocker(ctx, a.keeper)
	return nil
}
