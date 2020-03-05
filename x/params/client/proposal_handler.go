package client

import (
	govclient "github.com/Dipper-Protocol/x/gov/client"
	"github.com/Dipper-Protocol/x/params/client/cli"
	"github.com/Dipper-Protocol/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
