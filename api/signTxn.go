package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"

	"github.com/KushalP47/CSE542-Blockchain-Project/blockchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// type Tx struct {
// 	To    common.Address
// 	Value uint64
// 	Nonce uint64
// }

// type SignedTx struct {
// 	To      common.Address
// 	Value   uint64
// 	Nonce   uint64
// 	V, R, S *big.Int // signature values
// }

func SignTxn(w http.ResponseWriter, r *http.Request) {
	// Sign transaction
	var req struct {
		To    common.Address `json:"to"`
		Value uint64         `json:"value"`
		Nonce uint64         `json:"nonce"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := blockchain.Txn{
		To:    req.To,
		Value: req.Value,
		Nonce: req.Nonce,
	}

	// hash of txn
	h := Hash(&tx)
	privatekey, err := crypto.HexToECDSA("4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356")
	if err != nil {
		panic(err)
	}

	sig, err := crypto.Sign(h[:], privatekey)
	if err != nil {
		panic(err)
	}

	R, S, V := decodeSignature(sig)
	signedTx := blockchain.SignedTx{
		To:    tx.To,
		Value: tx.Value,
		Nonce: tx.Nonce,
		V:     V,
		R:     R,
		S:     S,
	}
	encodedSignedTxn, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(encodedSignedTxn))
	// Send signed transaction
	json.NewEncoder(w).Encode(hex.EncodeToString(encodedSignedTxn))
}

// HashSigned returns the tx hash
func HashSigned(tx *blockchain.SignedTx) common.Hash {
	return rlpHash(tx)
}

func Hash(tx *blockchain.Txn) common.Hash {
	return rlpHash([]interface{}{
		tx.To,
		tx.Value,
		tx.Nonce,
	})
}

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])

	return h
}

// hasherPool holds LegacyKeccak256 hashers for rlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// decodeSignature decodes the signature into v, r, and s values
func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}
