package api

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Credentials struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
	PublicKey  ecdsa.PublicKey
}

var myCredentials Credentials

func GenerateAddress(w http.ResponseWriter, r *http.Request) {
	// Generate random private key and address
	randomPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	randomAddress := crypto.PubkeyToAddress(randomPrivateKey.PublicKey)

	myCredentials = Credentials{
		PrivateKey: randomPrivateKey,
		Address:    randomAddress,
		PublicKey:  randomPrivateKey.PublicKey,
	}

	w.Write([]byte("Address generated successfully"))

}
