package service

import (
	"context"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rhizomplatform/plateaus/internal/polygon"
	"github.com/rhizomplatform/plateaus/internal/polygon/crypto"
	"github.com/rhizomplatform/plateaus/x/validation/client/rpc"
	"github.com/rhizomplatform/plateaus/x/validation/types"
)

type PlateausValidator struct {
	dialer polygon.Dialer
}

func New(dialer polygon.Dialer) PlateausValidator {
	return PlateausValidator{
		dialer: dialer,
	}
}

func (pv PlateausValidator) ConfirmValidator(ctx sdk.Context, externalKey string) (bool, error) {
	return ConfirmValidator(ctx.Context(), ctx.ChainID(), pv.dialer, externalKey)
}

func ConfirmValidator(
	ctx context.Context,
	plateausChainID string,
	dialer polygon.Dialer,
	pathExternalKey string,
) (bool, error) {
	fromAddress, privateKeyPolygon, err := crypto.AddressFromPrivateKey(pathExternalKey)

	if err != nil {
		return false, errors.New("invalid private key")
	}

	remoteRPC := types.GetRemoteRPC(plateausChainID)
	plateausValidatorTokenAddr := types.GetTokenAddr(plateausChainID)
	clientPolygon, err := dialer(remoteRPC)

	defer clientPolygon.Close()

	if err != nil {
		return false, errors.New(fmt.Sprintf("was not able to Dial: %s", err))
	}

	remoteChainId, err := clientPolygon.ChainID(ctx)

	if err != nil {
		return false, errors.New(fmt.Sprintf("was not able to ChainID: %s", err))
	}

	addressPlateausValidatorPolygon := common.HexToAddress(plateausValidatorTokenAddr)
	plateausValidatorContract, err := types.NewPlateausValidator(addressPlateausValidatorPolygon, clientPolygon)

	if err != nil {
		return false, errors.New(fmt.Sprintf("was not able to create NewPlateausValidator: %s", err))
	}

	rpcClient := rpc.New(clientPolygon, remoteChainId, plateausValidatorContract, *fromAddress, privateKeyPolygon)
	balance, err := rpcClient.BalanceOf()

	if err != nil {
		return false, errors.New(fmt.Sprintf("was not able to BalanceOf: %s", err))
	}

	return balance > 0, nil
}
