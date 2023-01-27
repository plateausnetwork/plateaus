package validation

import (
	"github.com/rhizomplatform/plateaus/x/validation/keeper"
	"github.com/rhizomplatform/plateaus/x/validation/types"
	"github.com/spf13/cast"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker sets the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	valAddr, err := sdk.ValAddressFromBech32(cast.ToString(k.ModuleOpts[types.ValidatorAddr]))

	if err != nil {
		k.Logger(ctx).Error("error on validator_address config", err)
		return
	}

	k.CheckValidator(ctx, valAddr)
}
