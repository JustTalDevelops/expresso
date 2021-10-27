package packet

import (
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/text"
)

// Disconnect is sent to the client when it's connection is closed.
type Disconnect struct {
	// Reason is the reason for closing the connection.
	Reason text.Text
}

// ID ...
func (*Disconnect) ID() int32 {
	return 0x1A
}

// Marshal ...
func (pk *Disconnect) Marshal(w *protocol.Writer) {
	w.Text(&pk.Reason)
}

// Unmarshal ...
func (pk *Disconnect) Unmarshal(r *protocol.Reader) {
	r.Text(&pk.Reason)
}
