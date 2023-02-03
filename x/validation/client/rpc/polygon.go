package rpc

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rhizomplatform/plateaus/x/validation/types"
	"math/big"
)

type Client struct {
	rpcClient   *ethclient.Client
	chainID     *big.Int
	contract    *types.PlateausValidator
	fromAddress common.Address
	privateKey  *ecdsa.PrivateKey
}

func New(rpc *ethclient.Client, chainId *big.Int, contract *types.PlateausValidator, fromAddress common.Address, privateKey *ecdsa.PrivateKey) *Client {
	return &Client{
		rpcClient:   rpc,
		chainID:     chainId,
		contract:    contract,
		fromAddress: fromAddress,
		privateKey:  privateKey,
	}
}

func (c Client) BalanceOf() (int64, error) {
	balance, err := c.contract.BalanceOf(nil, c.fromAddress)

	if err != nil {
		fmt.Printf("could not get balance of: %s", err)

		return 0, err
	}

	return balance.Int64(), nil
}
