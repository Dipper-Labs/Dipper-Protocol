package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/client/utils"
	govutils "github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov/client/utils"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov/types"
	"github.com/Dipper-Labs/Dipper-Protocol/client"
	"github.com/Dipper-Labs/Dipper-Protocol/client/context"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	"github.com/Dipper-Labs/Dipper-Protocol/version"
)

// Proposal flags
const (
	FlagTitle        = "title"
	FlagDescription  = "description"
	flagProposalType = "type"
	FlagDeposit      = "deposit"
	flagVoter        = "voter"
	flagDepositor    = "depositor"
	flagStatus       = "status"
	flagNumLimit     = "limit"
	FlagProposal     = "proposal"
)

type proposal struct {
	Title       string
	Description string
	Type        string
	Deposit     string
}

// ProposalFlags defines the core required fields of a proposal. It is used to
// verify that these values are not provided in conjunction with a JSON proposal
// file.
var ProposalFlags = []string{
	FlagTitle,
	FlagDescription,
	flagProposalType,
	FlagDeposit,
}

// GetTxCmd returns the transaction commands for this module
// governance ModuleClient is slightly different from other ModuleClients in that
// it contains a slice of "proposal" child commands. These commands are respective
// to proposal type handlers that are implemented in other modules but are mounted
// under the governance CLI (eg. parameter change proposals).
func GetTxCmd(storeKey string, cdc *codec.Codec, pcmds []*cobra.Command) *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Governance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmdSubmitProp := GetCmdSubmitProposal(cdc)
	for _, pcmd := range pcmds {
		cmdSubmitProp.AddCommand(client.PostCommands(pcmd)[0])
	}

	cmdSubmitProp.AddCommand(client.PostCommands(GetCmdSubmitSoftwareUpgradeProposal(cdc))[0])

	govTxCmd.AddCommand(client.PostCommands(
		GetCmdDeposit(cdc),
		GetCmdVote(cdc),
		cmdSubmitProp,
	)...)

	return govTxCmd
}

func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit a proposal along with an initial deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal along with an initial deposit.
Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Example:
$ %s tx gov submit-proposal --proposal="path/to/proposal.json" --from mykey

Where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10test"
}

Which is equivalent to:

$ %s tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10test" --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := parseSubmitProposalFlags()
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(proposal.Deposit)
			if err != nil {
				return err
			}

			content := types.ContentFromProposalType(proposal.Title, proposal.Description, proposal.Type)

			msg := types.NewMsgSubmitProposal(content, amount, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "title of proposal")
	cmd.Flags().String(FlagDescription, "", "description of proposal")
	cmd.Flags().String(flagProposalType, "", "proposalType of proposal, types: text/parameter_change/software_upgrade")
	cmd.Flags().String(FlagDeposit, "", "deposit of proposal")
	cmd.Flags().String(FlagProposal, "", "proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

func GetCmdSubmitSoftwareUpgradeProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "software-upgrade [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a software upgrade proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a software upgrade proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal software-upgrade <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
    "title":"testnet-v1.1.0 upgrade",
    "description":"upgrade for smart contract",
    "type":"SoftwareUpgrade",
    "deposit":{
        "denom":"pdip",
        "amount":"1000000"
    },
    "version":1,
    "software":"https://github.com/Dipper-Labs/Dipper-Protocol/releases/tag/testnet-v1.1.0",
    "switch_height":100000,
    "threshold":"90.000000000000000000"
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			var proposalJSON SoftwareUpgradeProposalJSON
			err = json.Unmarshal(contents, &proposalJSON)
			if err != nil {
				return err
			}

			var proposal types.SoftwareUpgradeProposal

			proposal.Title = proposalJSON.Title
			proposal.Description = proposalJSON.Description
			proposal.Version = proposalJSON.Version
			proposal.Software = proposalJSON.Software
			proposal.SwitchHeight = proposalJSON.SwitchHeight
			proposal.Threshold = proposalJSON.Threshold

			if err = proposal.ValidateBasic(); err != nil {
				return err
			}

			msg := types.NewMsgSubmitProposal(proposal, sdk.NewCoins(proposalJSON.Deposit), cliCtx.FromAddress)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [proposal-id] [deposit]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens for an active proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a deposit for an active proposal. You can
find the proposal-id by running "%s query gov proposals".

Example:
$ %s tx gov deposit 1 10stake --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			from := cliCtx.GetFromAddress()

			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(from, proposalID, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%s query gov proposals".


Example:
$ %s tx gov vote 1 yes --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			from := cliCtx.GetFromAddress()

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// Find out which vote option user chose
			byteVoteOption, err := types.VoteOptionFromString(govutils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := types.NewMsgVote(from, proposalID, byteVoteOption)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// DONTCOVER
