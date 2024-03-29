package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	service      service.PlateausValidator
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak distrtypes.AccountKeeper, blockedAddrs map[string]bool,
	moduleOpts map[string]interface{}, service service.PlateausValidator,
) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace,
		authKeeper:   ak,
		blockedAddrs: blockedAddrs,
		ModuleOpts:   moduleOpts,
		service:      service,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// CheckValidator check if validator is allowed to receive rewards
func (k Keeper) CheckValidator(ctx sdk.Context) {
	k.Logger(ctx).
		With("hash", ctx.HeaderHash().String()).
		Info("starting check validator permission")

	ok, err := k.service.ConfirmValidator(ctx, cast.ToString(k.ModuleOpts[types.FlagExternalKeyPathKey]))

	if err != nil || !ok {
		k.Logger(ctx).
			With("hash", ctx.HeaderHash().String()).
			Error("could not check validator", err)

		panic(err)
	}

	k.Logger(ctx).Info("validator was checked")

	return
}
