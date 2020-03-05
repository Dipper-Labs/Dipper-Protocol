package cli

import (
	"fmt"

	"github.com/Dipper-Protocol/client"
	"github.com/Dipper-Protocol/client/context"
	"github.com/Dipper-Protocol/codec"
	"github.com/Dipper-Protocol/x/dipperBank/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	dipperBankQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the dipperBank module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dipperBankQueryCmd.AddCommand(client.GetCommands(
		GetCmdOraclePrice(storeKey, cdc),
		GetCmdNetValue(storeKey, cdc),
		GetCmdBorrowBalance(storeKey, cdc),
		GetCmdBorrowValue(storeKey, cdc),
		GetCmdBorrowValueEstimate(storeKey, cdc),
		GetCmdSupplyBalance(storeKey, cdc),
		GetCmdSupplyValue(storeKey, cdc),
	)...)
	return dipperBankQueryCmd
}


// GetCmdOraclePrice queries a list of token price
func GetCmdOraclePrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracleprice [symbol]",
		Short: "Get oracle price",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			//addr := args[0]
			symbol := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oracleprice/%s", queryRoute, symbol), nil)
			if err != nil {
				fmt.Printf("could not get oracle price - %s \n", symbol)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdNetValue queries netValue of someone
func GetCmdNetValue(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "netvalue [addr]",
		Short: "net value of addr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr := args[0]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/netvalue/%s", queryRoute, addr), nil)
			if err != nil {
				fmt.Printf("could not get netValueOf - %s\n", addr)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdBorrowBalance queries a balance who has borrowed
func GetCmdBorrowBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "borrowbalance [addr] [symbol]",
		Short: "borrow balance of addr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]
			addr := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/borrowbalance/%s/%s", queryRoute, symbol, addr), nil)
			if err != nil {
				fmt.Printf("could not get borrow %s for %s \n", symbol, addr)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdBorrowValue queries value who has borrowed
func GetCmdBorrowValue(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "borrowValue [symbol] [addr]",
		Short: "borrow value of addr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]
			addr := args[1]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/borrowvalue/%s/%s", queryRoute, symbol, addr), nil)
			if err != nil {
				fmt.Printf("could not get borrowValueOf -symbol %s for %s\n", symbol, addr)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdBorrowValueEstimate queries value who can borrow
func GetCmdBorrowValueEstimate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "borrowvalueestimate [amount] [symbol]",
		Short: "borrow value estimate",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			amount := args[0]
			symbol := args[1]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/borrowvalueestimate/%s/%s", queryRoute, amount, symbol), nil)
			if err != nil {
				fmt.Printf("could not get borrowValueOf -symbol %s for %s\n", amount, symbol)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdSupplyBalance queries total supply balance
func GetCmdSupplyBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "supplybalance [symbol] [addr]",
		Short: "supply balance of",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]
			addr := args[1]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/supplybalance/%s/%s", queryRoute, symbol, addr), nil)
			if err != nil {
				fmt.Printf("could not get supplyBalanceOf -symbol %s for %s\n", symbol, addr)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdSupplyValue queries total value that every supplies
func GetCmdSupplyValue(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "supplyvalue [symbol] [addr]",
		Short: "supply value of addr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]
			addr := args[1]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/supplyvalue/%s/%s", queryRoute, symbol, addr), nil)
			if err != nil {
				fmt.Printf("could not get supplyValueOf -symbol %s for %s\n", symbol, addr)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}



