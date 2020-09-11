package cli

import (
	"github.com/spf13/cobra"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/types"
	"github.com/Dipper-Labs/Dipper-Protocol/client"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auth transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetMultiSignCommand(cdc),
		GetSignCommand(cdc),
	)
	return txCmd
}
