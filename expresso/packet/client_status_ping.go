package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// ClientStatusPing is sent by the client to the server to request a ServerStatusPong.
type ClientStatusPing struct {
	// Payload a system-dependent time value which is counted in milliseconds.
	Payload int64
}

// ID ...
func (*ClientStatusPing) ID() int32 {
	return 0x01
}

// Marshal ...
func (pk *ClientStatusPing) Marshal(w *protocol.Writer) {
	w.Int64(&pk.Payload)
}

// Unmarshal ...
func (pk *ClientStatusPing) Unmarshal(r *protocol.Reader) {
	r.Int64(&pk.Payload)
}
