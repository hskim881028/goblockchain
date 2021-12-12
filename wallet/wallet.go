package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/hskim881028/goblockchain/utility"
)

const (
	walletFileName string = "goblockchain.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	return !os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utility.HandleError(err)
	return privateKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utility.HandleError(err)
	err = os.WriteFile(walletFileName, bytes, 0644)
	utility.HandleError(err)
}

func restoreKey() *ecdsa.PrivateKey {
	bytes, err := os.ReadFile(walletFileName)
	utility.HandleError(err)
	privateKey, err := x509.ParseECPrivateKey(bytes)
	utility.HandleError(err)
	return privateKey
}

func encodeBigInts(a, b []byte) string {
	sum := append(a, b...)
	return fmt.Sprintf("%x", sum)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Sign(payload string, w *wallet) string {
	paylaodAsBytes, err := hex.DecodeString(payload)
	utility.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, paylaodAsBytes)
	utility.HandleError(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func Verify(signature, address, payload string) bool {
	x, y, err := restoreBigInts(address)
	utility.HandleError(err)
	publicKey := ecdsa.PublicKey{elliptic.P256(), x, y}

	payloadAsBytes, err := hex.DecodeString(payload)
	utility.HandleError(err)

	r, s, err := restoreBigInts(signature)
	utility.HandleError(err)

	ok := ecdsa.Verify(&publicKey, payloadAsBytes, r, s)
	return ok
}
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
