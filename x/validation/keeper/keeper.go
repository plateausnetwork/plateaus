package keeper

import (
	"github.com/rhizomplatform/plateaus/x/validation/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
		Info("Starting check validator permission")

	if k.isInitialized(ctx, valAddr) {
		return
	}

	//TODO: get permission external request
	value := true
	k.SetValidator(ctx, valAddr, value)

	return
}

func (k Keeper) isInitialized(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
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
