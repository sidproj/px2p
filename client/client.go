package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"px2p/protocol"
)

func handleMessage(p *protocol.Peer, msg protocol.Message) {
	switch msg.Type {
	case "PONG":
		fmt.Println("Server: PONG")
	case "CHAT":
		fmt.Println("Server: ", string(msg.Data))
	case "ERROR":
		fmt.Println("Server error: ", string(msg.Data))
	default:
		fmt.Println("Unknown message: ", msg.Type)
	}
}

func StartClient(conn net.Conn) {

	peer := protocol.NewPeer(conn, handleMessage)

	go peer.ReadLoop()
	go peer.WriteLoop()

	reader := bufio.NewReader(os.Stdin)

	peer.Send <- protocol.Message{Type: "PING"}

	for {
		fmt.Print("Client: ")
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		peer.Send <- protocol.Message{
			Type: "CHAT",
			Data: []byte(data),
		}

		if data == "EXIT" {
			peer.Done <- struct{}{}
		}
	}
}

func Connect() {

	conn, err := net.Dial("tcp", "localhost:4040")

	fmt.Println("Connected to the server.")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer conn.Close()

	err = protocol.ClientPerformHandshake(conn, "peer-1", "v1")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Handshake successful!")

	StartClient(conn)

}
