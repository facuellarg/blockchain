package main

import (
	"blockchain/blockchain"
	"blockchain/wallet"
	"fmt"
)

func main() {
	// b := NewBlock(3, "initHash")
	// fmt.Println(b.Hash())
	// b.Print()
	// myAddress := "my_address"
	// bc := blockchain.NewBlockChain(myAddress)
	// bc.Print()
	// bc.AddTransaction("A", "B", 0.5)
	// bc.Mining()
	// bc.Print()
	// bc.AddTransaction("C", "D", 13)
	// bc.AddTransaction("X", "Y", 12)
	// bc.Mining()
	// bc.Print()
	// fmt.Printf("bc.CalculateTotalAmount(myAddress): %v\n", bc.CalculateTotalAmount(myAddress))
	// fmt.Printf("bc.CalculateTotalAmount(\"A\"): %v\n", bc.CalculateTotalAmount("A"))
	// fmt.Printf("bc.CalculateTotalAmount(\"B\"): %v\n", bc.CalculateTotalAmount("B"))
	// fmt.Printf("bc.CalculateTotalAmount(\"D\"): %v\n", bc.CalculateTotalAmount("D"))

	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	transaction := wallet.NewTransaction(
		walletA.PrivateKey(),
		walletA.PublicKey(),
		walletA.BlockchainAddress(),
		walletB.BlockchainAddress(),
		1.0,
	)
	bc := blockchain.NewBlockChain(walletM.BlockchainAddress(), 5000)
	isAdded := bc.AddTransaction(
		walletA.BlockchainAddress(),
		walletB.BlockchainAddress(),
		1.0,
		walletA.PublicKey(),
		transaction.GenerateSignature(),
	)
	bc.Mining()
	bc.Print()
	fmt.Printf("isAdded: %v\n", isAdded)
	fmt.Printf("bc.CalculateTotalAmount(walletA.BlockchainAddress()): %v\n", bc.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("bc.CalculateTotalAmount(walletB.BlockchainAddress()): %v\n", bc.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("bc.CalculateTotalAmount(walletM.BlockchainAddress()): %v\n", bc.CalculateTotalAmount(walletM.BlockchainAddress()))
}
