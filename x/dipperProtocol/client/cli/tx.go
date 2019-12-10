package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/Dipper-Protocol/x/dipperProtocol/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	dipperProtocolTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "dipperProtocol transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	dipperProtocolTxCmd.AddCommand(client.PostCommands(
		GetCmdBuyName(cdc),
		GetCmdSetName(cdc),
		GetCmdDeleteName(cdc),
		GetCmdBankBorrow(cdc),
		GetCmdBankDeposit(cdc),
		GetCmdBankWithdraw(cdc),
		GetCmdBankRepay(cdc),
		GetCmdSetOraclePrice(cdc),
	)...)

	return dipperProtocolTxCmd
}

// GetCmdBuyName is the CLI command for sending a BuyName transaction
func GetCmdBuyName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy-name [name] [amount]",
		Short: "bid for existing name or claim new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBuyName(args[0], coins, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSetName is the CLI command for sending a SetName transaction
func GetCmdSetName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-name [name] [value]",
		Short: "set the value associated with a name that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// if err := cliCtx.EnsureAccountExists(); err != nil {
			// 	return err
			// }

			msg := types.NewMsgSetName(args[0], args[1], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDeleteName is the CLI command for sending a DeleteName transaction
func GetCmdDeleteName(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-name [name]",
		Short: "delete the name that you own along with it's associated fields",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgDeleteName(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}


// GetCmdDeleteName is the CLI command for sending a DeleteName transaction
func GetCmdBankBorrow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bank-borrow [amount] [symbol]",
		Short: "borrow from the pool, if you have deposit",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBankBorrow(coins, args[1], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBankRepay is the CLI command for sending a BankRepay transaction
func GetCmdBankRepay(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bank-repay [amount] [symbol]",
		Short: "repay to the pool, if you have borrowed",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgBankRepay(coins, args[1], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBankDeposit is the CLI command for sending a BankDeposit transaction
func GetCmdBankDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bank-deposit [name] [symbol]",
		Short: "deposit to bank, if you have money that the bank supports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[0])
			if err != nil{
				return err
			}

			msg := types.NewMsgBankDeposit(coins, args[1], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDeleteName is the CLI command for sending a DeleteName transaction
func GetCmdBankWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bank-withdraw [name] [symbol]",
		Short: "withdraw from bank, if you have money deposit in the bank",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[0])
			if err != nil{
				return err
			}

			msg := types.NewMsgBankWithdraw(coins, args[1], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSetOraclePrice is the CLI command for sending a SetOraclePrice transaction
func GetCmdSetOraclePrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-oracleprice [name] [symbol] [amount]",
		Short: "set the oracle price",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			//coins, err := sdk.ParseCoins(args[0])
			//if err != nil{
			//	return err
			//}

			msg := types.NewMsgSetOraclePrice(args[0], args[1], args[2], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}