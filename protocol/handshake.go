package protocol

import (
	"encoding/json"
	"fmt"
	"net"
)

type Handshake struct {
	Type    string `json:"type"`
	PeerID  string `json:"peer_id"`
	Version string `json:"version"`
}

func ClientPerformHandshake(conn net.Conn, peerID string, version string) error {

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	hello := Handshake{
		Type:    "HELLO",
		PeerID:  peerID,
		Version: version,
	}

	if err := enc.Encode(hello); err != nil {
		return err
	}

	var resp Handshake

	if err := dec.Decode(&resp); err != nil {
		return err
	}

	if resp.Type != "WELCOME" {
		return fmt.Errorf("Handshake rejected!")
	}
	return nil
}

func ServerPerformHandshake(conn net.Conn) (string, error) {

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	var hello Handshake

	if err := dec.Decode(&hello); err != nil {
		return "", err
	}

	if hello.Type != "HELLO" {
		return "", fmt.Errorf("Invalid handshake")
	}

	if hello.Version != "v1" {
		return "", fmt.Errorf("Unsoported version")
	}

	resp := Handshake{
		Type:    "WELCOME",
		PeerID:  "server",
		Version: "v1",
	}

	if err := enc.Encode(resp); err != nil {
		return "", err
	}

	return hello.PeerID, nil
}
