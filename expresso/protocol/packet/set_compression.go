package packet

import (
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

// SetCompression updates the compression threshold on the client to whatever threshold the
// server requests. This is usually done right after encryption.
type SetCompression struct {
	// Threshold is the compression threshold to be used.
	Threshold int32
}

// ID ...
func (*SetCompression) ID() int32 {
	return 0x03
}

// Marshal ...
func (pk *SetCompression) Marshal(w *protocol.Writer) {
	w.Varint32(&pk.Threshold)
}

// Unmarshal ...
func (pk *SetCompression) Unmarshal(r *protocol.Reader) {
	r.Varint32(&pk.Threshold)
}
