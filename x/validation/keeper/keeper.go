package keeper

import (
	"fmt"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/rhizomplatform/plateaus/x/validation/service"
	"github.com/rhizomplatform/plateaus/x/validation/types"
	"github.com/spf13/cast"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the distribution store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace
	authKeeper distrtypes.AccountKeeper

	blockedAddrs map[string]bool
	ModuleOpts   map[string]interface{}
	memStore     sdk.KVStore
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak distrtypes.AccountKeeper, memStore sdk.KVStore, blockedAddrs map[string]bool, moduleOpts map[string]interface{},
) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace,
		authKeeper:   ak,
		memStore:     memStore,
		blockedAddrs: blockedAddrs,
		ModuleOpts:   moduleOpts,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// CheckValidator check if validator is allowed to receive rewards
func (k Keeper) CheckValidator(ctx sdk.Context, valAddr sdk.ValAddress) {
	k.Logger(ctx).
		With("hash", ctx.HeaderHash().String()).
		With("validator", valAddr.String()).
		Info("starting check validator permission")

	externalAdd := cast.ToString(k.ModuleOpts[types.ExternalAddrKey])
	validations, err := service.GetValidations(valAddr, externalAdd)

	if err != nil {
		k.Logger(ctx).
			With("hash", ctx.HeaderHash().String()).
			With("validator", valAddr.String()).
			Error("could not get validations", err)

		return
	}

	for receivedAddr, isAble := range validations {
		accAddr, err := sdk.AccAddressFromBech32(receivedAddr)

		if err != nil {
			k.Logger(ctx).With("received-addr", receivedAddr).Error("received acc address was not able to create an AccAddress")
			continue
		}

		k.Logger(ctx).With("received-addr", receivedAddr).Error("setting validator validation")

		val := sdk.ValAddress(accAddr.Bytes())
		k.SetValidator(ctx, val, isAble)
	}

	k.Logger(ctx).Info("validator was checked")

	return
}

func (k Keeper) createTx(ctx sdk.Context, cdc codec.ProtoCodecMarshaler, msg sdk.Msg) {
	txCfg := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	txBuilder := txCfg.NewTxBuilder()

	if err := txBuilder.SetMsgs(msg); err != nil {
		panic(err)
	}

	bz, err := txCfg.TxEncoder()(txBuilder.GetTx())

	if err != nil {
		panic(err)
	}

	ctxClient := sdkclient.Context{}.
		WithChainID(ctx.ChainID()).
		WithCodec(cdc).
		WithBroadcastMode(flags.BroadcastSync).
		WithViper("PLATEAUS")

	ctxClient, _ = config.ReadFromClientConfig(ctxClient)

	// Broadcast the transaction
	res, err := ctxClient.BroadcastTx(bz)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.Code)
	fmt.Println(res)

}

func (k Keeper) SetValidator(ctx sdk.Context, valAddr sdk.ValAddress, value bool) {
	bValue := []byte("false")

	if value {
		bValue = []byte("true")
	}

	k.memStore.Set(types.GetValidatorValidationRewardsKey(valAddr), bValue)
}

func (k Keeper) HasPermission(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	value := k.memStore.Get(types.GetValidatorValidationRewardsKey(valAddr))

	if value == nil {
		return false
	}

	sValue := string(value)

	if sValue == "true" {
		return true
	}

	return false
}
