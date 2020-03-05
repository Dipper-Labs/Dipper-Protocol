package main

import (
	"encoding/json"
	app "github.com/Dipper-Protocol/app"
	"github.com/Dipper-Protocol/baseapp"
	"github.com/Dipper-Protocol/client"
	"github.com/Dipper-Protocol/store"
	"github.com/spf13/viper"
	"io"

	"github.com/Dipper-Protocol/server"
	"github.com/Dipper-Protocol/x/genaccounts"
	genaccscli "github.com/Dipper-Protocol/x/genaccounts/client/cli"
	"github.com/Dipper-Protocol/x/staking"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/Dipper-Protocol/types"
	genutilcli "github.com/Dipper-Protocol/x/genutil/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagOverwrite    = "overwrite"
	flagMinGasPrices = "minimum-gas-prices"
)

// dipd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "dipd",
		Short:             "dipperProtocol App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		client.NewCompletionCmd(rootCmd, true),
		replayCmd(),
		client.LineBreak,
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "DIP", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
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

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

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

// AddGenesisAccountCmd allows users to add accounts to the genesis file
//func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "add-genesis-account [address] [coins[,coins]]",
//		Short: "Adds an account to the genesis file",
//		Args:  cobra.ExactArgs(2),
//		Long: strings.TrimSpace(`
//Adds accounts to the genesis file so that you can start a chain with coins in the CLI:
//
//$ dipd add-genesis-account cosmos1tse7r2fadvlrrgau3pa0ss7cqh55wrv6y9alwh 1000STAKE,1000nametoken
//`),
//		RunE: func(_ *cobra.Command, args []string) error {
//			addr, err := sdk.AccAddressFromBech32(args[0])
//			if err != nil {
//				return err
//			}
//			coins, err := sdk.ParseCoins(args[1])
//			if err != nil {
//				return err
//			}
//			coins.Sort()
//
//			var genDoc tmtypes.GenesisDoc
//			config := ctx.Config
//			genFile := config.GenesisFile()
//			if !common.FileExists(genFile) {
//				return fmt.Errorf("%s does not exist, run `gaiad init` first", genFile)
//			}
//			genContents, err := ioutil.ReadFile(genFile)
//			if err != nil {
//			}
//
//			if err = cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
//				return err
//			}
//
//			var appState app.GenesisState
//			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
//				return err
//			}
//			//TODO temporary case
//			authGenState := auth.GetGenesisStateFromAppState(cdc, appState)
//			if authGenState.Accounts.Contains(addr) {
//				return fmt.Errorf("cannot add account at existing address %s", addr)
//			}
//
//			acc := auth.NewBaseAccountWithAddress(addr)
//			acc.Coins = coins
//
//			authGenState.Accounts = append(authGenState.Accounts, &acc)
//			authGenState.Accounts = auth.SanitizeGenesisAccounts(authGenState.Accounts)
//			authGenStateBz, err := cdc.MarshalJSON(authGenState)
//			if err != nil {
//				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
//			}
//
//			appState[auth.ModuleName] = authGenStateBz
//			appStateJSON, err := cdc.MarshalJSON(appState)
//			if err != nil {
//				return fmt.Errorf("failed to marshal application genesis state: %w", err)
//			}
//
//			genDoc.AppState = appStateJSON
//
//			return genutil.ExportGenesisFile(&genDoc, genFile)
//			//return gaiaInit.ExportGenesisFile(genFile, genDoc.ChainID, genDoc.Validators, appStateJSON)
//		},
//	}
//	return cmd
//}

// SimpleAppGenTx returns a simple GenTx command that makes the node a valdiator from the start
//func SimpleAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (
//	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
//
//	addr, secret, err := server.GenerateCoinKey()
//	if err != nil {
//		return
//	}
//
//	bz, err := cdc.MarshalJSON(struct {
//		Addr sdk.AccAddress `json:"addr"`
//	}{addr})
//	if err != nil {
//		return
//	}
//
//	appGenTx = json.RawMessage(bz)
//
//	bz, err = cdc.MarshalJSON(map[string]string{"secret": secret})
//	if err != nil {
//		return
//	}
//
//	cliPrint = json.RawMessage(bz)
//
//	validator = tmtypes.GenesisValidator{
//		PubKey: pk,
//		Power:  10,
//	}
//
//	return
//}