package client

import (
	govclient "github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov/client"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params/client/cli"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params/client/rest"
)

// ProposalHandler - param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
