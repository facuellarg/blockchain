package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {

	wallet := new(Wallet)
	wallet.privateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	wallet.publicKey = &wallet.privateKey.PublicKey
	// 2 perform sha 256 hashin the public key
	h2 := sha256.New()
	h2.Write(wallet.publicKey.X.Bytes())
	h2.Write(wallet.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	//3 perform RIPEMD-160 hashing on the above result
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	//4 add version byte in front of the above hash
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	//5 perform sha 256 with the above result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	//6 perform hash 256 with the above hash
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7 take the first 4 bytes of the second hash for checksum
	chsum := digest6[:4]
	//8 add the above 4 digest to the end of extended version of RIPMED-160
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])
	//9 convert the result from a byte to string into base58
	address := base58.Encode(dc8)

	wallet.blockchainAddress = address
	return wallet
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Address    string `json:"address"`
	}{
		PrivateKey: w.PrivateKeyStr(),
		PublicKey:  w.PublicKeyStr(),
		Address:    w.BlockchainAddress(),
	})
}
