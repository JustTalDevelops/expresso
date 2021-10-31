package packet

// State represents the state the connection is on. These split up the packets that can be accessed
// to a specific selection based on the state.
type State struct {
	// packetCollection is the packet collection the state is linked to.
	*packetCollection
	// state is the actual state value.
	state
}

// StateHandshaking represents the handshaking state.
func StateHandshaking() State {
	return State{state: stateHandshaking, packetCollection: handshakingCollection}
}

// StateStatus represents the status state.
func StateStatus() State {
	return State{state: stateStatus, packetCollection: statusCollection}
}

// StateLogin represents the login state.
func StateLogin() State {
	return State{state: stateLogin, packetCollection: loginCollection}
}

// StatePlay represents the play state.
func StatePlay() State {
	return State{state: statePlay, packetCollection: playCollection}
}

// Packet finds a packet in the state. based on the target direction and ID.
func (s State) Packet(direction Direction, id int32) Packet {
	packetMap := s.packetMap(direction)
	if packetMap[id] == nil {
		return nil
	}

	return packetMap[id]()
}

// packetMap returns the packet map for the direction.
func (s State) packetMap(direction Direction) map[int32]func() Packet {
	switch direction {
	case DirectionServer():
		return s.clientBoundPackets
	case DirectionClient():
		return s.serverBoundPackets
	}
	panic("should never happen")
}

type state byte

const (
	stateHandshaking state = iota
	stateStatus
	stateLogin
	statePlay
)
