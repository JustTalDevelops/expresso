package packet

import (
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

// ServerKeepAlive is a packet sent by the server to the client usually every two seconds
// to keep the connection alive.
type ServerKeepAlive struct {
	// PingID is the time in Unix milliseconds that the server sent the packet.
	PingID int64
}

// ID ...
func (*ServerKeepAlive) ID() int32 {
	return 0x21
}

// Marshal ...
func (pk *ServerKeepAlive) Marshal(w *protocol.Writer) {
	w.Int64(&pk.PingID)
}

// Unmarshal ...
func (pk *ServerKeepAlive) Unmarshal(r *protocol.Reader) {
	r.Int64(&pk.PingID)
}
