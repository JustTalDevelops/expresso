package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// ClientKeepAlive is a packet sent by the client to the server usually every two seconds
// to keep the connection alive.
type ClientKeepAlive struct {
	// PingID is the time in Unix milliseconds that the client sent the packet.
	PingID int64
}

// ID ...
func (*ClientKeepAlive) ID() int32 {
	return 0x0F
}

// Marshal ...
func (pk *ClientKeepAlive) Marshal(w *protocol.Writer) {
	w.Int64(&pk.PingID)
}

// Unmarshal ...
func (pk *ClientKeepAlive) Unmarshal(r *protocol.Reader) {
	r.Int64(&pk.PingID)
}
