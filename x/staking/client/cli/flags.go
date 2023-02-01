package cli

import (
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	flag "github.com/spf13/pflag"
)

const (
	FlagExternalAddress = "external-addr"
)

// FlagSetExternalAddress Returns the flagset for External Address related operations.
func FlagSetExternalAddress() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagExternalAddress, "", "The external address to check your node license validator")
	return fs
}

func flagSetDescriptionCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(cli.FlagMoniker, "", "The validator's name")
	fs.String(cli.FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fs.String(cli.FlagWebsite, "", "The validator's (optional) website")
	fs.String(cli.FlagSecurityContact, "", "The validator's (optional) security contact email")
	fs.String(cli.FlagDetails, "", "The validator's (optional) details")

	return fs
}
