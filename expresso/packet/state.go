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
	return State{
		packetCollection: handshakingCollection,
		state:            stateHandshaking,
	}
}

// StateStatus represents the status state.
func StateStatus() State {
	return State{
		packetCollection: statusCollection,
		state:            stateStatus,
	}
}

// StateLogin represents the login state.
func StateLogin() State {
	return State{
		packetCollection: loginCollection,
		state:            stateLogin,
	}
}

// StatePlay represents the play state.
func StatePlay() State {
	return State{
		packetCollection: playCollection,
		state:            statePlay,
	}
}

// Packet finds a packet in the state. based on the target bound and ID.
func (s State) Packet(bound Bound, id int32) Packet {
	if bound == BoundClient() {
		return s.clientBoundPackets[id]()
	}
	return s.serverBoundPackets[id]()
}

type state byte

const (
	stateHandshaking state = iota
	stateStatus
	stateLogin
	statePlay
)
