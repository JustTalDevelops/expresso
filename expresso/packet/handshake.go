package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// Handshake is sent when the client initially tries to join the server. It is the first packet sent and contains
// information specific to the player.
type Handshake struct {
	// Protocol is the protocol version of the player. The player is disconnected if the protocol is
	// incompatible with the protocol of the server.
	Protocol int32
	// Address is the address the player used to connect to the server.
	Address string
	// Port is the port of the server that the player used to connect with.
	Port int16
	// NextState is either one for status, or two for login.
	NextState int32
}

// ID ...
func (*Handshake) ID() int32 {
	return 0x00
}

// Marshal ...
func (pk *Handshake) Marshal(w *protocol.Writer) {
	w.Varint32(&pk.Protocol)
	w.String(&pk.Address)
	w.Int16(&pk.Port)
	w.Varint32(&pk.NextState)
}

// Unmarshal ...
func (pk *Handshake) Unmarshal(r *protocol.Reader) {
	r.Varint32(&pk.Protocol)
	r.String(&pk.Address)
	r.Int16(&pk.Port)
	r.Varint32(&pk.NextState)
}
