package blockchain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	senderBlockchainAddress   string
	recipientBlockchainAddres string
	value                     float64
}

func NewTransaction(
	sender,
	recipient string,
	value float64) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("sender :\t%s\n", t.senderBlockchainAddress)
	fmt.Printf("recipient :\t%s\n", t.recipientBlockchainAddres)
	fmt.Printf("value :\t%.4f\n", t.value)

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
