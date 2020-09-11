package cli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/types"
	"github.com/Dipper-Labs/Dipper-Protocol/client"
	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	"github.com/Dipper-Labs/Dipper-Protocol/hexutil"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/Dipper-Labs/Dipper-Protocol/version"
)

var ZeroAmount = sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	vmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for vm",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	vmQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(cdc),
		GetCmdQueryDBState(cdc),
		GetCmdQueryCode(cdc),
		GetCmdGetStorage(cdc),
		GetCmdGetLogs(cdc),
		GetCmdQueryCreateFee(cdc),
		GetCmdQueryCallFee(cdc),
		GetCmdQueryCall(cdc),
	)...)
	return vmQueryCmd
}

func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current vm parameters",
		Long: strings.TrimSpace(fmt.Sprintf(`Query values set as vm parameters.
Example:
$ %s query vm params`, version.ClientName)),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Println("fail")
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

func GetCmdQueryDBState(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state [--all] [--show_code]",
		Short: "Query the current vm whole state",
		Long: strings.TrimSpace(fmt.Sprintf(`Query values set as vm state.
Example:
$ %s query vm state [--all] [--show_code]`, version.ClientName)),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			var params types.QueryStateParams

			t := cmd.Flag(flagAll)
			if t != nil {
				params.ContractOnly = !t.Changed
			}

			t = cmd.Flag(flagShowCode)
			if t != nil {
				params.ShowCode = t.Changed
			}

			pd, err := json.Marshal(params)

			fmt.Println(string(pd))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryState)
			bz, _, err := cliCtx.QueryWithData(route, pd)
			if err != nil {
				return err
			}

			var out bytes.Buffer
			err = json.Indent(&out, bz, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(out.String())

			return nil
		},
	}

	cmd.Flags().Bool(flagAll, false, "Show all account object, including non contract account object")
	cmd.Flags().Bool(flagShowCode, false, "Show contract code")

	return cmd
}

func GetCmdQueryCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "code",
		Short: "Querying commands for Contract Code",
		Long: strings.TrimSpace(fmt.Sprintf(`Query Contract Code by Account Address.
Example:
$ %s query vm code [address]`, version.ClientName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			_, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/vm/%s/%s", types.QueryCode, args[0])
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("no code found with address %s", args[0])
			}

			dst := make([]byte, 2*len(res))
			hex.Encode(dst, res)

			fmt.Println(string(dst))

			return nil
		},
	}
}

func GetCmdGetStorage(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "storage [account] [key]",
		Short: "Querying storage for an account at a given key",
		Long: strings.TrimSpace(fmt.Sprintf(`Query Contract Code by Account Address.
Example:
$ %s query vm code [address]`, version.ClientName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/vm/%s/%s/%s", types.QueryStorage, addr, args[1])
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}

			var out types.QueryStorageResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdGetLogs(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "logs [txhash]",
		Short: "Querying logs by txHash",
		Long: strings.TrimSpace(fmt.Sprintf(`Query logs by txHash.
Example:
$ %s query vm logs [txHash]`, version.ClientName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.Query(
				fmt.Sprintf("custom/vm/logs/%s", args[0]))
			if err != nil {
				return err
			}

			var out types.QueryLogsResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdQueryCreateFee(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "feecreate [code_file]",
		Short: "Querying fee to deploy contract",
		Long: strings.TrimSpace(fmt.Sprintf(`Querying fee to deploy contract.
Example:
$ %s query vm feecreate [code_file] [from_accaddr]`, version.ClientName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			code, _ := CodeFromFile(args[0])
			msg := types.NewMsgContractQuery(from, nil, code, ZeroAmount)
			data, err := cliCtx.Codec.MarshalJSON(msg)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/vm/%s", types.EstimateGas), data)
			if err != nil {
				return err
			}

			var out types.SimulationResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdQueryCallFee(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "feecall [from] [to] [method] [args] [amount] [abi_file]",
		Short: "Querying fee to call contract",
		Long: strings.TrimSpace(fmt.Sprintf(`Querying fee to call contract.
Example:
$ %s query vm feecall dip1mfztsv6eq5rhtaz2l6jjp3yup3q80agsqra9qe dip1rk47h83x4nz4745d63dtnpl8uwsramfgz8snr5 balanceOf 0000000000000000000000000000000000000000000000000000000000000001 0pdip ./demo.abi`, version.ClientName)),
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			abiObj, err := AbiFromFile(args[5])
			if err != nil {
				return err
			}

			argsBin, err := hex.DecodeString(args[3])
			if err != nil {
				return err
			}

			method := args[2]
			m, exist := abiObj.Methods[method]
			var payload []byte
			if exist {
				if len(m.Inputs) != len(argsBin)/32 {
					return errors.New("args count dismatch")
				}

				readyArgs, err := m.Inputs.UnpackValues(argsBin)
				if err != nil {
					return err
				}

				payload, err = abiObj.Pack(method, readyArgs...)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("method %s not exist", method)
			}

			msg := types.NewMsgContractQuery(fromAddr, toAddr, payload, ZeroAmount)
			data, err := cliCtx.Codec.MarshalJSON(msg)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/vm/%s", types.EstimateGas), data)
			if err != nil {
				return err
			}

			var out types.SimulationResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdQueryCall(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call [from] [to] [method] [args] [amount] [abi_file]",
		Short: "Querying fee to call contract",
		Long: strings.TrimSpace(fmt.Sprintf(`call contract for query, don't create a transaction.
Example:
$ %s query vm call dip1mfztsv6eq5rhtaz2l6jjp3yup3q80agsqra9qe dip1rk47h83x4nz4745d63dtnpl8uwsramfgz8snr5 balanceOf ./demo.abi --amount=0pdip --args="arg1 arg2"`, version.ClientName)),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			argList := viper.GetStringSlice(flagArgs)
			payload, m, err := GenPayload(args[3], args[2], argList)
			if err != nil {
				return err
			}

			msg := types.NewMsgContractQuery(fromAddr, toAddr, payload, ZeroAmount)
			data, err := cliCtx.Codec.MarshalJSON(msg)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/vm/%s", types.QueryCall), data)
			if err != nil {
				return err
			}

			var out types.SimulationResult
			cdc.MustUnmarshalJSON(res, &out)

			d, err := hexutil.Decode(out.Res)
			if err != nil {
				return cliCtx.PrintOutput(out)
			}

			var result types.VMQueryResult
			result.Gas = out.Gas
			result.Result, err = m.Outputs.UnpackValues(d)
			if err != nil {
				return err
			}

			for i := 0; i < len(m.Outputs); i++ {
				if m.Outputs[i].Type.String() == "address" {
					var addr sdk.AccAddress
					if ethAddr, ok := result.Result[i].(common.Address); ok {
						addr = append(addr[:], ethAddr.Bytes()...)
						result.Result[i] = addr.String()
					}
				}
			}

			fmt.Println(result)

			return nil

		},
	}

	cmd.Flags().String(flagArgs, "", "contract method arg list (e.g. --args='arg1 arg2 arg3')(default \"\")")
	cmd.Flags().String(flagAmount, "0pdip", "amount of coins to send (e.g. --amount=100pdip)")

	return cmd
}
