package main

import (
	"blockchain/blockchain"
	"blockchain/utils"
	"blockchain/wallet"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockChainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockChain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = blockchain.NewBlockChain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
	}
	return bc
}

func HelloWorld(ctx echo.Context) error {
	return ctx.String(200, "Hello")
}

func (bcs *BlockchainServer) MakeTransaction(ctx echo.Context) error {
	transactionRequest := &blockchain.TransactionRequest{}
	if err := ctx.Bind(transactionRequest); err != nil {
		log.Error(err)
		return err
	}
	bc := bcs.GetBlockChain()
	publicKey := utils.PublicKeyFromString(transactionRequest.SenderPublicKey)
	// privateKey := utils.PrivateKeyFromString(
	// transactionRequest.SenderPublicKey,
	// publicKey,
	// )
	if bc.CreateTransaction(
		transactionRequest.SenderBlockchainAddress,
		transactionRequest.RecipientBlockchainAddress,
		transactionRequest.Value,
		publicKey,
		transactionRequest.Signature,
	) {
		return ctx.NoContent(http.StatusCreated)
	}
	return echo.ErrBadRequest
}

func (bcs *BlockchainServer) GetTransactions(ctx echo.Context) error {
	bc := bcs.GetBlockChain()
	transactions := bc.TransactionPool()
	return ctx.JSON(http.StatusOK, echo.Map{
		"transactions": transactions,
		"length":       len(transactions),
	})
}

func (bcs *BlockchainServer) StartMining(ctx echo.Context) error {
	bc := bcs.GetBlockChain()
	bc.StartMining()
	return ctx.NoContent(http.StatusOK)
}

func (bcs *BlockchainServer) GetChain(ctx echo.Context) error {
	bc := bcs.GetBlockChain()
	m, err := json.Marshal(bc)
	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, string(m[:]))
}

func (bcs *BlockchainServer) Amount(ctx echo.Context) error {
	address := ctx.QueryParam("address")
	if address == "" {
		return echo.ErrBadRequest
	}
	bc := bcs.GetBlockChain()
	amount := bc.CalculateTotalAmount(address)
	return ctx.JSON(http.StatusOK, echo.Map{
		"amount": amount,
	})
}
func (bcs *BlockchainServer) Run() error {
	server := echo.New()
	server.Use(middleware.Logger())
	server.GET("/", bcs.GetChain)
	server.POST("/transactions", bcs.MakeTransaction)
	server.GET("/mine", bcs.StartMining)
	server.GET("/amount", bcs.Amount)
	server.GET("/transactions", bcs.GetTransactions)
	return server.Start(fmt.Sprintf(":%d", bcs.Port()))
}
