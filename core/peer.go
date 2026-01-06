package core

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	Conn   net.Conn
	PeerId string
	enc    *json.Encoder
	dec    *json.Decoder

	Send      chan Message
	Done      chan struct{}
	handler   MessageHandler
	closeOnce sync.Once
}

func NewPeer(conn net.Conn, peerID string, handler MessageHandler) *Peer {
	return &Peer{
		Conn:    conn,
		PeerId:  peerID,
		enc:     json.NewEncoder(conn),
		dec:     json.NewDecoder(conn),
		Send:    make(chan Message, 16),
		Done:    make(chan struct{}),
		handler: handler,
	}
}

func (p *Peer) ReadLoop(peerEvent chan PeerEvent) {
	defer p.Close()

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

func (p *Peer) SendMessage(msg Message) bool {
	select {
	case <-p.Done:
		return false
	case p.Send <- msg:
		return true
	}
}

func (p *Peer) Close() {
	p.closeOnce.Do(func() {
		close(p.Done)
		p.Conn.Close()
	})
}
