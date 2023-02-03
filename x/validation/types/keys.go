package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "validation"

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// QuerierRoute is the querier route for distribution
	QuerierRoute = ModuleName

	// FlagExternalKeyPathKey is the key for config
	FlagExternalKeyPathKey = "external-key-path"

	// EveryTimeKey is the key for config
	EveryTimeKey = "every-time"
)

// Keys for distribution store
// Items are stored with the following key: values
//
// - 0x00<proposalID_Bytes>: FeePol
//
// - 0x01: sdk.ConsAddress
//
// - 0x02<valAddrLen (1 Byte)><valAddr_Bytes>: ValidatorValidationRewards
var (
	ValidatorValidationRewardsPrefix = []byte{0x02} // key for outstanding rewards
)

// GetValidatorValidationRewardsAddress creates an address from a validator's outstanding rewards key.
func GetValidatorValidationRewardsAddress(key []byte) (valAddr sdk.ValAddress) {
	// key is in the format:
	// 0x02<valAddrLen (1 Byte)><valAddr_Bytes>

	// Remove prefix and address length.
	kv.AssertKeyAtLeastLength(key, 3)
	addr := key[2:]
	kv.AssertKeyLength(addr, int(key[1]))

	return sdk.ValAddress(addr)
}

// GetValidatorValidationRewardsKey creates the validation rewards key for a validator.
func GetValidatorValidationRewardsKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorValidationRewardsPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}
