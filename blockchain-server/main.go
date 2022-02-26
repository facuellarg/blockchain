package main

import "flag"

var (
	port *uint
)

func init() {
	port = flag.Uint("p", 5000, "TCP port number for blockchain server")

}

func main() {
	flag.Parse()
	server := NewBlockChainServer(uint16(*port))
	server.Run()
}
