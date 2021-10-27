package packet

import (
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

// EncryptionResponse is sent back by the client to agree to the encryption. Every future
// packet will use the agreed upon encryption.
type EncryptionResponse struct {
	// SharedSecret is a shared secret which was encrypted with the server's public key.
	SharedSecret []byte
	// VerifyToken is the same verify token value encrypted with the server's public key.
	VerifyToken []byte
}

// ID ...
func (*EncryptionResponse) ID() int32 {
	return 0x01
}

// Marshal ...
func (pk *EncryptionResponse) Marshal(w *protocol.Writer) {
	w.ByteSlice(&pk.SharedSecret)
	w.ByteSlice(&pk.VerifyToken)
}

// Unmarshal ...
func (pk *EncryptionResponse) Unmarshal(r *protocol.Reader) {
	r.ByteSlice(&pk.SharedSecret)
	r.ByteSlice(&pk.VerifyToken)
}
