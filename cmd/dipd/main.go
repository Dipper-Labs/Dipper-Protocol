package main

import (
	"encoding/json"
	"io"

	"github.com/Dipper-Labs/Dipper-Protocol/baseapp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Dipper-Labs/Dipper-Protocol/app"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/genaccounts"
	genaccscli "github.com/Dipper-Labs/Dipper-Protocol/app/v0/genaccounts/client/cli"
	genutilcli "github.com/Dipper-Labs/Dipper-Protocol/app/v0/genutil/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/guardian"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking"
	"github.com/Dipper-Labs/Dipper-Protocol/client"
	"github.com/Dipper-Labs/Dipper-Protocol/server"
	"github.com/Dipper-Labs/Dipper-Protocol/store"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	flagMinGasPrices   = "minimum-gas-prices"
	flagInvCheckPeriod = "inv-check-period"
)

var invCheckPeriod uint

func main() {
	cdc := app.MakeLatestCodec()

	config := sdk.GetConfig()
	config.Seal()

	ctx := server.NewDefaultContext()

	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "dipd",
		Short:             "dip Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(guardian.AddGenesisGuardianCmd(ctx, cdc, app.DefaultNodeHome))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(replayCmd())
	rootCmd.AddCommand(client.LineBreak)
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "DIP", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	minGasPrices := viper.GetString(flagMinGasPrices)
	return app.NewDIPApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))), baseapp.SetMinGasPrices(minGasPrices),
	)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		dipApp := app.NewDIPApp(logger, db, traceStore, false, uint(1))
		err := dipApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return dipApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	dipApp := app.NewDIPApp(logger, db, traceStore, true, uint(1))
	return dipApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
