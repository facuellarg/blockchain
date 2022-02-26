package main

import (
	"flag"
)

var (
	port    *uint
	gateway *string
)

func main() {
	port = flag.Uint("p", 8080, "tcp port for wallet server")
	gateway = flag.String("gw", "http://127.0.0.1:5000", "tcp port for wallet server")
	flag.Parse()
	ws := NewWalletServer(uint16(*port), *gateway)
	ws.Run()
}
