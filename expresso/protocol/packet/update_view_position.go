package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// UpdateViewPosition is sent by the server to update the client's view position.
type UpdateViewPosition struct {
	// X, Z are the coordinates of the new chunk to be viewed.
	X, Z int32
}

// ID ...
func (*UpdateViewPosition) ID() int32 {
	return 0x49
}

// Marshal ...
func (pk *UpdateViewPosition) Marshal(w *protocol.Writer) {
	w.Varint32(&pk.X)
	w.Varint32(&pk.Z)
}

// Unmarshal ...
func (pk *UpdateViewPosition) Unmarshal(r *protocol.Reader) {
	r.Varint32(&pk.X)
	r.Varint32(&pk.Z)
}
