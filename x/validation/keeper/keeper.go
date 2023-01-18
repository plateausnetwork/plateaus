package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	plateaustypes "github.com/rhizomplatform/plateaus/types"
	"github.com/rhizomplatform/plateaus/x/validation/types"
	"github.com/tendermint/tendermint/libs/log"
	"net/http"
	"time"
)

// Keeper of the distribution store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace
	authKeeper distrtypes.AccountKeeper

	blockedAddrs map[string]bool
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak distrtypes.AccountKeeper, blockedAddrs map[string]bool,
) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace,
		authKeeper:   ak,
		blockedAddrs: blockedAddrs,
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
		With("validator", valAddr.Bytes()).
		Info("starting check validator permission")

	if k.isInitialized(ctx, valAddr) {
		return
	}

	valAccAddr, _ := plateaustypes.GetPlateausAddressFromBech32(valAddr.String())

	//TODO: get permission external request
	//TODO: techinical debt: we need to execute this validation onchain
	c := &http.Client{
		Timeout: time.Millisecond * 5000,
	}
	r, err := c.Get(fmt.Sprintf("https://vyd0yyst26.execute-api.us-east-1.amazonaws.com/development/nodes/%s", valAccAddr))

	if err != nil {
		k.Logger(ctx).
			With("hash", ctx.HeaderHash().String()).
			With("validator", valAddr.Bytes()).
			Error("could not check validator permission")
		return
	}

	defer r.Body.Close()

	value := r.StatusCode == http.StatusNoContent

	k.Logger(ctx).Info(fmt.Sprintf("validator was checked: %t", value))

	k.SetValidator(ctx, valAddr, value)

	return
}

func (k Keeper) isInitialized(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorValidationRewardsPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		k.Logger(ctx).With("validator", key).With("permission", value).Info(fmt.Sprintf("validators was initialized"))
	}

	return store.Has(types.GetValidatorValidationRewardsKey(valAddr))
}

func (k Keeper) SetValidator(ctx sdk.Context, valAddr sdk.ValAddress, value bool) {
	bValue := []byte("false")

	if value {
		bValue = []byte("true")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorValidationRewardsKey(valAddr), bValue)
}

func (k Keeper) HasPermission(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetValidatorValidationRewardsKey(valAddr))

	if value == nil {
		return false
	}

	sValue := string(value)

	if sValue == "true" {
		return true
	}

	return false
}
