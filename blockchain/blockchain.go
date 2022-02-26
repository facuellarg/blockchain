package blockchain

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	DifficultyMining       = 3
	MiningSender           = "TheBlockchain"
	BlockchainMiningReward = 1.0
	MiningTimeSec          = 10
)

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	address         string
	port            uint16
	mux             sync.Mutex
}

type TransactionRequest struct {
	SenderPublicKey            string           `json:"sender_public_key"`
	SenderBlockchainAddress    string           `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string           `json:"recipient_blockchain_address"`
	Value                      float64          `json:"value,string"`
	Signature                  *utils.Signature `json:"signature"`
}

func NewBlockChain(address string, port uint16) *Blockchain {
	b := &Block{}
	blockChain := new(Blockchain)
	blockChain.CreateBlock(0, b.Hash())
	blockChain.address = address
	blockChain.port = port
	return blockChain
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, block)
	bc.transactionPool = []*Transaction{}
	return block
}

func (bc *Blockchain) CreateTransaction(
	sender,
	recipient string,
	value float64,
	senderPublicKey *ecdsa.PublicKey,
	signature *utils.Signature,
) bool {

	isTransacted := bc.AddTransaction(
		sender,
		recipient,
		value,
		senderPublicKey,
		signature,
	)

	return isTransacted

}

func (bc *Blockchain) AddTransaction(
	sender,
	recipient string,
	value float64,
	senderPublicKey *ecdsa.PublicKey,
	signature *utils.Signature,
) bool {

	t := &Transaction{sender, recipient, value}
	if sender == MiningSender {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, signature, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: Not enough balance in your wallet")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	log.Println("ERROR: signature verification fail")
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey,
	signature *utils.Signature,
	transaction *Transaction,
) bool {
	m, err := json.Marshal(transaction)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], signature.R, signature.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, len(bc.transactionPool))
	for i, v := range bc.transactionPool {
		transactions[i] = NewTransaction(v.senderBlockchainAddress, v.recipientBlockchainAddres, v.value)
	}
	return transactions

}

func (bc *Blockchain) ValidProof(nonce int, prevHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, prevHash, transactions}
	guessBlockHash := fmt.Sprintf("%x", guessBlock.Hash())
	return guessBlockHash[:difficulty] == zeros
}

//ProofOfWork
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	prevHash := bc.LastBlock().previousHash
	nonce := 0
	for !bc.ValidProof(nonce, prevHash, transactions, DifficultyMining) {
		nonce++
	}
	return nonce
}

//LastBlock
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

//Print
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain:%d\t%s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Println(strings.Repeat("*", 50))
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	if len(bc.TransactionPool()) == 0 {
		return false
	}
	bc.AddTransaction(MiningSender, bc.address, BlockchainMiningReward, nil, nil)
	nonce := bc.ProofOfWork()
	prevHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, prevHash)
	log.Println("action=mining,status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	time.AfterFunc(MiningTimeSec*time.Second, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(address string) float64 {
	totalAmount := 0.0
	for _, block := range bc.chain {
		for _, transaction := range block.transactions {
			if transaction.recipientBlockchainAddres == address {
				totalAmount += transaction.value
			}
			if transaction.senderBlockchainAddress == address {
				totalAmount -= transaction.value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}
