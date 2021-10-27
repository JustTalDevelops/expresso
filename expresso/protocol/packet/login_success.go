package packet

import (
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

// LoginSuccess is sent by the server to the client to notify it that it's login attempt was successful.
// It is the last packet of the login sequence, and is the indicator to switch to the play state.
type LoginSuccess struct {
	// UUID is the UUID of the player logging in.
	UUID uuid.UUID
	// Username is the username of the player logging in.
	Username string
}

// ID ...
func (*LoginSuccess) ID() int32 {
	return 0x02
}

// Marshal ...
func (pk *LoginSuccess) Marshal(w *protocol.Writer) {
	w.UUID(&pk.UUID)
	w.String(&pk.Username)
}

// Unmarshal ...
func (pk *LoginSuccess) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.UUID)
	r.String(&pk.Username)
}
