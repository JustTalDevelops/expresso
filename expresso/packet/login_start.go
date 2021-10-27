package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// LoginStart is sent by the client to request starting a login.
type LoginStart struct {
	// Username is the username of the player.
	Username string
}

// ID ...
func (*LoginStart) ID() int32 {
	return 0x00
}

// Marshal ...
func (pk *LoginStart) Marshal(w *protocol.Writer) {
	w.String(&pk.Username)
}

// Unmarshal ...
func (pk *LoginStart) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Username)
}
