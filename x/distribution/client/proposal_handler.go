package client

import (
	"github.com/Dipper-Protocol/x/distribution/client/cli"
	"github.com/Dipper-Protocol/x/distribution/client/rest"
	govclient "github.com/Dipper-Protocol/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
