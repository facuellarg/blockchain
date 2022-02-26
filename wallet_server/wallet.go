package main

import (
	"blockchain/blockchain"
	"blockchain/templates"
	"blockchain/utils"
	"blockchain/wallet"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WalletServer struct {
	port    uint16
	gateway string
}

type TransactionRequest struct {
	SenderPrivateKey           string  `json:"sender_private_key"`
	SenderPublicKey            string  `json:"sender_public_key"`
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float64 `json:"value,string"`
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "hello", "")
}

func (ws *WalletServer) Wallet(ctx echo.Context) error {
	myWallet := wallet.NewWallet()
	if err := ctx.Bind(myWallet); err != nil {
		log.Println(err)
		return err
	}
	return ctx.JSON(http.StatusOK, myWallet)
}

func (ws *WalletServer) Path(path string) string {
	return fmt.Sprintf("%s/%s", ws.Gateway(), path)
}

func (ws *WalletServer) GetTotalAmount(ctx echo.Context) error {
	address := ctx.QueryParam("address")
	if address == "" {
		return echo.ErrBadRequest
	}
	client := resty.New()
	resp, err := client.R().
		SetQueryParam("address", address).
		Get(ws.Path("amount"))
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadGateway
	}

	amount := struct {
		Amount float64 `json:"amount"`
	}{}
	if err := json.Unmarshal(resp.Body(), &amount); err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadGateway
	}
	return ctx.JSON(http.StatusOK, amount)
}

func (ws *WalletServer) CreateTransaction(ctx echo.Context) error {
	transactionRequest := &TransactionRequest{}
	if err := ctx.Bind(transactionRequest); err != nil {
		log.Println(err)
		return err
	}
	publicKey := utils.PublicKeyFromString(transactionRequest.SenderPublicKey)
	privateKey := utils.PrivateKeyFromString(transactionRequest.SenderPrivateKey, publicKey)

	transaction := wallet.NewTransaction(
		privateKey,
		publicKey,
		transactionRequest.SenderBlockchainAddress,
		transactionRequest.RecipientBlockchainAddress,
		transactionRequest.Value,
	)
	signature := transaction.GenerateSignature()
	bcTransaction := blockchain.TransactionRequest{
		SenderPublicKey:            transactionRequest.SenderPublicKey,
		SenderBlockchainAddress:    transactionRequest.SenderBlockchainAddress,
		RecipientBlockchainAddress: transactionRequest.RecipientBlockchainAddress,
		Signature:                  signature,
		Value:                      transactionRequest.Value,
	}
	data, _ := json.Marshal(bcTransaction)
	resp, err := http.Post(
		ws.Path("transactions"),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	return ctx.NoContent(resp.StatusCode)
	// return ctx.JSON(http.StatusOK, transactionRequest)
}

func (ws *WalletServer) Run() error {
	t := templates.NewTemplate(template.Must(template.ParseGlob("./../views/wallet/*.html")))
	server := echo.New()
	server.Renderer = t
	server.Use(middleware.Logger())
	server.GET("/", ws.Index)
	server.POST("/wallet", ws.Wallet)
	server.GET("/amount", ws.GetTotalAmount)
	server.POST("/transaction", ws.CreateTransaction)
	return server.Start(fmt.Sprintf(":%d", ws.Port()))
}
