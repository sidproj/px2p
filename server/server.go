package server

import (
	"fmt"
	"log"
	"net"
	"px2p/protocol"
	"time"
)

type Handshake struct {
	Type    string `json:"type"`
	PeerID  string `json:"peer_id"`
	Version string `json:"version"`
}

func handleMessage(p *protocol.Peer, msg protocol.Message) {
	switch msg.Type {
	case "PING":
		fmt.Println("Client: PING")
		p.Send <- protocol.Message{Type: "PONG"}
	case "CHAT":
		fmt.Println("Client: ", string(msg.Data))
		p.Send <- protocol.Message{Type: "CHAT", Data: msg.Data}
	case "ERROR":
		fmt.Println("Client error: ", string(msg.Data))
	default:
		fmt.Println("Unknown message: ", msg.Type)
	}
}

func Comunicate(conn net.Conn) {

	peer := protocol.NewPeer(conn, handleMessage)

	go peer.ReadLoop()
	go peer.WriteLoop()

	<-peer.Done

}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	// set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	fmt.Println("New connection: ", conn.RemoteAddr(), conn.LocalAddr())

	peerID, err := protocol.ServerPerformHandshake(conn)
	if err != nil {
		fmt.Println("Handshake failed: ", err)
		return
	}

	conn.SetDeadline(time.Time{})
	fmt.Println("Handshake successful with: ", peerID)
	Comunicate(conn)
}

func StartServer() {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		fmt.Println("Error while creating a server, trying another port")
		ln, err = net.Listen("tcp", ":4050")
		if err != nil {
			log.Fatal("Error while creating a server")
		}
	}

	fmt.Println("Started listening on port:", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
