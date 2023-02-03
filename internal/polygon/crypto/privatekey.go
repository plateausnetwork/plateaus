package crypto

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	cryptoethereum "github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"strings"
)

func AddressFromPrivateKey(pKey string) (*common.Address, *ecdsa.PrivateKey, error) {
	key, err := getExternalKey(pKey)

	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	privateKey, err := cryptoethereum.HexToECDSA(key)

	if err != nil {
		log.Printf("could not get address crypto.HexToECDSA: %s", err)
		return nil, nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		err := errors.New(fmt.Sprintf("could not get ecdsa.PublicKey: %t", ok))
		log.Println(err)
		return nil, nil, err
	}

	fromAddress := cryptoethereum.PubkeyToAddress(*publicKeyECDSA)

	return &fromAddress, privateKey, nil
}

func getExternalKey(path string) (string, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return "", errors.New("path not found or invalid content")
	}

	return strings.TrimRight(string(content), "\n"), nil
}
