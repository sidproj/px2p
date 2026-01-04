package main

import (
	"px2p/client"
	"px2p/server"
)

func main() {
	go server.StartServer()
	go client.Connect()
	for {

	}
}
