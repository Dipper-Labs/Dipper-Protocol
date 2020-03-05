package vm

import (
	"encoding/json"
	cli2 "github.com/Dipper-Protocol/x/vm/client/cli"
	rest2 "github.com/Dipper-Protocol/x/vm/client/rest"
	"github.com/Dipper-Protocol/x/vm/types"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Protocol/client/context"
	"github.com/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(value json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(value, &data); err != nil {
		return err
	}

	return ValidateGenesis(data)
}

func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest2.RegisterRoutes(ctx, rtr)
}

func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli2.GetQueryCmd(types.StoreKey, cdc)
}

var _ module.AppModuleBasic = AppModuleBasic{}

type AppModule struct {
	AppModuleBasic
	k Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{k: keeper}
}

func (a AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	a.k.SetParams(ctx, genesisState.Params)

	return nil
}

func (a AppModule) ExportGenesis(sdk.Context) json.RawMessage {
	return nil
}

func (a AppModule) RegisterInvariants(sdk.InvariantRegistry) {
	panic("implement me")
}

func (a AppModule) Route() string {
	return RouterKey
}

func (a AppModule) NewHandler() sdk.Handler {
	return NewHandler(a.k)
}

func (a AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (a AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(a.k)
}

func (a AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {
	// TODO
}

func (a AppModule) EndBlock(ctx sdk.Context, end abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, a.k)
}
