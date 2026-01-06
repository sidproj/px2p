package server

import (
	"fmt"
	"log"
	"net"
	"px2p/core"
	"px2p/protocol"
	"time"
)

type Handshake struct {
	Type    string `json:"type"`
	PeerID  string `json:"peer_id"`
	Version string `json:"version"`
}

func handleMessage(p *core.Peer, msg core.Message) {
	switch msg.Type {
	case "PING":
		fmt.Println("Client: PING")
		p.Send <- core.Message{Type: "PONG"}
	case "CHAT":
		fmt.Println("Client: ", string(msg.Data))
		p.Send <- core.Message{Type: "CHAT", Data: msg.Data}
	case "ERROR":
		fmt.Println("Client error: ", string(msg.Data))
	default:
		fmt.Println("Unknown message: ", msg.Type)
	}
}

func Comunicate(conn net.Conn, peer *core.Peer, peerEvent chan core.PeerEvent) {

	go peer.ReadLoop(peerEvent)
	go peer.WriteLoop()

	<-peer.Done

}

func handleConnection(conn net.Conn, peerEvent chan core.PeerEvent) {

	defer conn.Close()

	// ------------------- handshake ritual start -------------------
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	fmt.Println("New connection: ", conn.RemoteAddr(), conn.LocalAddr())

	peerID, err := protocol.ServerPerformHandshake(conn)
	if err != nil {
		fmt.Println("Handshake failed: ", err)
		return
	}
	conn.SetDeadline(time.Time{})
	fmt.Println("Handshake successful with: ", peerID)
	// ------------------- handshake ritual end -------------------

	peer := core.NewPeer(conn, peerID, handleMessage)
	peerInfo, err := core.NewPeerInfo(peer)

	if err != nil {
		fmt.Println("Unable to create peerInfo")
		return
	}

	peerEvent <- core.PeerEvent{
		Type: core.PEER_JOINED,
		Peer: peerInfo,
	}

	Comunicate(conn, peer, peerEvent)
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

	serverState := core.ServerState{
		PeerInfoList: core.PeerInfoList{
			Max:   3,
			Peers: make(map[string]core.PeerInfo),
		},
		Events: make(chan core.PeerEvent, 10),
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn, serverState.Events)
	}
}
