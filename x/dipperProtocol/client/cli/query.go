package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/Dipper-Protocol/x/dipperProtocol/internal/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	dipperProtocolQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the dipperProtocol module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dipperProtocolQueryCmd.AddCommand(client.GetCommands(
		GetCmdResolveName(storeKey, cdc),
		GetCmdWhois(storeKey, cdc),
		GetCmdNames(storeKey, cdc),
		GetCmdOraclePrice(storeKey, cdc),
		GetCmdNetValue(storeKey, cdc),
		GetCmdBorrowBalance(storeKey, cdc),
		GetCmdBorrowValue(storeKey, cdc),
		GetCmdBorrowValueEstimate(storeKey, cdc),
		GetCmdSupplyBalance(storeKey, cdc),
		GetCmdSupplyValue(storeKey, cdc),
	)...)
	return dipperProtocolQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [name]",
		Short: "resolve name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve name - %s \n", name)
				return nil
			}

			var out types.QueryResResolve
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdWhois queries information about a domain
func GetCmdWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "whois [name]",
		Short: "Query whois info of name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve whois - %s \n", name)
				return nil
			}

			var out types.Whois
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdNames queries a list of all names
func GetCmdNames(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "names",
		Short: "names",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/names", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get query names\n")
				return nil
			}

			var out types.QueryResNames
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdNames queries a list of all names
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

// GetCmdNames queries a list of all names
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

// GetCmdNames queries a list of all names
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

// GetCmdNames queries a list of all names
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

// GetCmdNames queries a list of all names
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

// GetCmdNames queries a list of all names
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

// GetCmdNames queries a list of all names
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



