package protocol

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Type string `json:"type"` // "CHAT", "PING", "PONG", "ERROR" , "STORE", "GET"
	Data []byte `json:"data,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type MessageHandler func(p *Peer, msg Message)

type Peer struct {
	conn net.Conn
	enc  *json.Encoder
	dec  *json.Decoder

	Send    chan Message
	Done    chan struct{}
	handler MessageHandler
}

func NewPeer(conn net.Conn, handler MessageHandler) *Peer {
	return &Peer{
		conn:    conn,
		enc:     json.NewEncoder(conn),
		dec:     json.NewDecoder(conn),
		Send:    make(chan Message, 16),
		Done:    make(chan struct{}),
		handler: handler,
	}
}

func (p *Peer) ReadLoop() {
	defer close(p.Done)

	for {
		var msg Message
		if err := p.dec.Decode(&msg); err != nil {
			fmt.Println("Read error:", err)
			return
		}

		p.handler(p, msg)
	}
}

func (p *Peer) WriteLoop() {
	for {
		select {
		case msg := <-p.Send:
			if err := p.enc.Encode(msg); err != nil {
				fmt.Println("Write error:", err)
				return
			}
		case <-p.Done:
			return
		}
	}
}
