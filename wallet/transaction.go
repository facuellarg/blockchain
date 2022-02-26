package wallet

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
)

type Transaction struct {
	senderPrivateKey          *ecdsa.PrivateKey
	senderPublicKey           *ecdsa.PublicKey
	senderBlockchainAddress   string
	recipientBlockchainAddres string
	value                     float64
}

func NewTransaction(
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	sender,
	recipient string,
	value float64,
) *Transaction {
	return &Transaction{
		privateKey,
		publicKey,
		sender,
		recipient,
		value,
	}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Value     float64 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddres,
		Value:     t.value,
	})

}
