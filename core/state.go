package core

type ServerState struct {
	PeerInfoList PeerInfoList
	Events       chan PeerEvent
}

func (s *ServerState) Run() {
	for event := range s.Events {
		switch event.Type {
		case PEER_JOINED:
			s.handlePeerJoin(event)
		}
	}
}

func (s *ServerState) handlePeerJoin(event PeerEvent) {
	if err := s.PeerInfoList.Add(event.Peer); err != nil {
		event.Reply <- PEER_REDIRECTED
		return
	}
	event.Reply <- PEER_ACCEPTED
}
