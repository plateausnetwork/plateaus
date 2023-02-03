package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rhizomplatform/plateaus/internal/polygon"
	"github.com/rhizomplatform/plateaus/x/validation/service"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// default values
var (
	DefaultTokens = sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
)

func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pathExternalKey, _ := cmd.Flags().GetString(FlagExternalKeyPath)
			ok, err := service.ConfirmValidator(context.Background(), clientCtx.ChainID, polygon.Dial, pathExternalKey)

			if err != nil {
				return err
			}

			if !ok {
				return errors.New("does not have permission")
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(cli.FlagSetPublicKey())
	cmd.Flags().AddFlagSet(cli.FlagSetAmount())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(cli.FlagSetCommissionCreate())
	cmd.Flags().AddFlagSet(cli.FlagSetMinSelfDelegation())

	cmd.Flags().String(cli.FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(cli.FlagNodeID, "", "The node's ID")

	cmd.Flags().AddFlagSet(FlagSetExternalAddress())

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(cli.FlagAmount)
	_ = cmd.MarkFlagRequired(cli.FlagPubKey)
	_ = cmd.MarkFlagRequired(cli.FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagExternalKeyPath)

	return cmd
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateValidator, error) {
	fAmount, _ := fs.GetString(cli.FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(cli.FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(cli.FlagMoniker)
	identity, _ := fs.GetString(cli.FlagIdentity)
	website, _ := fs.GetString(cli.FlagWebsite)
	security, _ := fs.GetString(cli.FlagSecurityContact)
	details, _ := fs.GetString(cli.FlagDetails)
	description := types.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	// get the initial validator commission parameters
	rateStr, _ := fs.GetString(cli.FlagCommissionRate)
	maxRateStr, _ := fs.GetString(cli.FlagCommissionMaxRate)
	maxChangeRateStr, _ := fs.GetString(cli.FlagCommissionMaxChangeRate)

	commissionRates, err := buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txf, nil, err
	}

	// get the initial validator min self delegation
	msbStr, _ := fs.GetString(cli.FlagMinSelfDelegation)

	minSelfDelegation, ok := sdk.NewIntFromString(msbStr)
	if !ok {
		return txf, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	msg, err := types.NewMsgCreateValidator(
		sdk.ValAddress(valAddr), pk, amount, description, commissionRates, minSelfDelegation,
	)
	if err != nil {
		return txf, nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(cli.FlagIP)
		nodeID, _ := fs.GetString(cli.FlagNodeID)

		if nodeID != "" && ip != "" {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txf, msg, nil
}
