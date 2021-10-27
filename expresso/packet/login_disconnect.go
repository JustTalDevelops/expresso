package packet

import (
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/text"
)

// LoginDisconnect is sent to the client when it's connection is closed during login.
type LoginDisconnect struct {
	// Reason is the reason for closing the connection.
	Reason text.Text
}

// ID ...
func (*LoginDisconnect) ID() int32 {
	return 0x00
}

// Marshal ...
func (pk *LoginDisconnect) Marshal(w *protocol.Writer) {
	w.Text(&pk.Reason)
}

// Unmarshal ...
func (pk *LoginDisconnect) Unmarshal(r *protocol.Reader) {
	r.Text(&pk.Reason)
}
