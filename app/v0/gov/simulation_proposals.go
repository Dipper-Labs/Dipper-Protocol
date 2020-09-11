package gov

// DONTCOVER

import (
	"math/rand"

	simappparams "github.com/Dipper-Labs/Dipper-Protocol/app/simapp/params"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov/types"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/simulation"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	simtypes "github.com/Dipper-Labs/Dipper-Protocol/types/simulation"
)

// OpWeightSubmitTextProposal app params key for text proposal
const OpWeightSubmitTextProposal = "op_weight_submit_text_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightMsgDeposit,
			simappparams.DefaultWeightTextProposal,
			SimulateTextProposalContent,
		),
	}
}

// SimulateTextProposalContent returns a random text proposal content.
func SimulateTextProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) simtypes.Content {
	return types.NewTextProposal(
		simtypes.RandStringOfLength(r, 140),
		simtypes.RandStringOfLength(r, 5000),
	)
}
