package types

import (
	"bytes"
	"fmt"

	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingProposalID uint64        `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Deposits           Deposits      `json:"deposits" yaml:"deposits"`
	Votes              Votes         `json:"votes" yaml:"votes"`
	Proposals          []Proposal    `json:"proposals" yaml:"proposals"`
	DepositParams      DepositParams `json:"deposit_params" yaml:"deposit_params"`
	VotingParams       VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams        TallyParams   `json:"tally_params" yaml:"tally_params"`
}

// NewGenesisState creates a new genesis state for the governance module
func NewGenesisState(startingProposalID uint64, dp DepositParams, vp VotingParams, tp TallyParams) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositParams:      dp,
		VotingParams:       vp,
		TallyParams:        tp,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(10)
	return GenesisState{
		StartingProposalID: 1,
		DepositParams: DepositParams{
			MinDeposit:       sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, minDepositTokens)},
			MaxDepositPeriod: DefaultDepositPeriod,
		},
		VotingParams: VotingParams{
			VotingPeriod: DefaultVotingPeriod,
		},
		TallyParams: TallyParams{
			Quorum:    sdk.NewDecWithPrec(334, 3),
			Threshold: sdk.NewDecWithPrec(5, 1),
			Veto:      sdk.NewDecWithPrec(334, 3),
		},
	}
}

// Equal- Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty - Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote threshold should be positive and less or equal to one, is %s", threshold.String())
	}

	veto := data.TallyParams.Veto
	if veto.IsNegative() || veto.GT(sdk.OneDec()) {
		return fmt.Errorf("governance vote veto threshold should be positive and less or equal to one, is %s", veto.String())
	}

	if !data.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("governance deposit amount must be a valid sdk.Coins amount, is %s", data.DepositParams.MinDeposit.String())
	}

	return nil
}
