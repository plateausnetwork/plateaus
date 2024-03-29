package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/rhizomplatform/plateaus/x/incentives/client/cli"
	"github.com/rhizomplatform/plateaus/x/incentives/client/rest"
)

var (
	RegisterIncentiveProposalHandler = govclient.NewProposalHandler(cli.NewRegisterIncentiveProposalCmd, rest.RegisterIncentiveProposalRESTHandler)
	CancelIncentiveProposalHandler   = govclient.NewProposalHandler(cli.NewCancelIncentiveProposalCmd, rest.CancelIncentiveProposalRequestRESTHandler)
)
