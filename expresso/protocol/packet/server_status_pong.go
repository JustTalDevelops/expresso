package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// ServerStatusPong is sent by the server to the client as a response to a ClientStatusPing.
type ServerStatusPong struct {
	// Payload is usually the same payload as the one in ClientStatusPing.
	Payload int64
}

// ID ...
func (*ServerStatusPong) ID() int32 {
	return 0x01
}

// Marshal ...
func (pk *ServerStatusPong) Marshal(w *protocol.Writer) {
	w.Int64(&pk.Payload)
}

// Unmarshal ...
func (pk *ServerStatusPong) Unmarshal(r *protocol.Reader) {
	r.Int64(&pk.Payload)
}
