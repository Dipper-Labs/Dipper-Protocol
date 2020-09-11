package client

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution/client/rest"
	govclient "github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov/client"
)

// ProposalHandler - param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
