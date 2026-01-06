package core

import (
	"fmt"
	"net"
	"time"
)

type PeerEventType int

type ServerEventResponseType int

const (
	PEER_JOINED PeerEventType = iota
	PEER_LEFT
	PEER_PING
)

const (
	PEER_ACCEPTED ServerEventResponseType = iota
	PEER_REDIRECTED
)

type PeerEvent struct {
	Type  PeerEventType
	Peer  PeerInfo
	Reply chan ServerEventResponseType
}

type Address struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type PeerInfo struct {
	PeerID   string    `json:"peer_id"`
	Addr     Address   `json:"address"`
	LastPing time.Time `json:"last_ping"`
	State    string    // "CONNECTED", "KNOWN", "DEAD"
}

type PeerInfoList struct {
	Peers map[string]PeerInfo
	Max   int
}

func (pl *PeerInfoList) Add(p PeerInfo) error {
	if pl.Max >= len(pl.Peers) {
		return fmt.Errorf("Peers connections are at max capacity")
	}

	pl.Peers[p.PeerID] = p
	return nil
}

func (pl *PeerInfoList) Remove(peerID string) {
	delete(pl.Peers, peerID)
}

func (pl *PeerInfoList) Snapshot() []PeerInfo {
	copy := make([]PeerInfo, 0, len(pl.Peers))
	for _, p := range pl.Peers {
		copy = append(copy, p)
	}
	return copy
}

func NewPeerInfo(peer *Peer) (PeerInfo, error) {
	tcpAddr, ok := peer.Conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return PeerInfo{}, fmt.Errorf("Invalid connection address")
	}

	ip := tcpAddr.IP.String()
	port := tcpAddr.Port

	return PeerInfo{
		PeerID: peer.PeerId,
		Addr: Address{
			IP:   ip,
			Port: port,
		},
		LastPing: time.Now(),
		State:    "KNOWN",
	}, nil
}
