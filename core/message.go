package core

type Message struct {
	Type string `json:"type"` // "CHAT", "PING", "PONG", "ERROR"
	// "REDIRECT"
	Data []byte `json:"data,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type MessageHandler func(p *Peer, msg Message)
