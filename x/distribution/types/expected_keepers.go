package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidationKeeper defines the expected interface needed to retrieve account balances.
type ValidationKeeper interface {
	CheckValidator(ctx sdk.Context, valAddr sdk.ValAddress)
	HasPermission(ctx sdk.Context, valAddr sdk.ValAddress) bool
}
