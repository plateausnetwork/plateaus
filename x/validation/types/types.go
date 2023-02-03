package types

import "github.com/rhizomplatform/plateaus/types"

// RemoteRPCMainnet remote addr for mainnet connection
const RemoteRPCMainnet = ""

// RemoteRPCTestnet remote addr for testnet connection
const RemoteRPCTestnet = "https://matic-mumbai.chainstacklabs.com"

// PlateausValidatorAddrMainnet token address for mainnet
const PlateausValidatorAddrMainnet = ""

// PlateausValidatorAddrTestnet token address for testnet
const PlateausValidatorAddrTestnet = "0x947e508c683314f5c99ddd88f17aca6ddf57f589"

// GetRemoteRPC to return the remote addr based on which network are running it
func GetRemoteRPC(chainID string) string {
	if types.IsMainnet(chainID) {
		return RemoteRPCMainnet
	}

	return RemoteRPCTestnet
}

// GetTokenAddr to return the remote addr based on which network are running it
func GetTokenAddr(chainID string) string {
	if types.IsMainnet(chainID) {
		return PlateausValidatorAddrMainnet
	}

	return PlateausValidatorAddrTestnet
}
