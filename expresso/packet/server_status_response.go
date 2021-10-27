package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// ServerStatusResponse is a response to the ClientStatusRequest packet, with the status sent as
// a string encoded in JSON.
type ServerStatusResponse struct {
	// Status is the JSON encoded status string.
	Status string
}

// ID ...
func (*ServerStatusResponse) ID() int32 {
	return 0x00
}

// Marshal ...
func (pk *ServerStatusResponse) Marshal(w *protocol.Writer) {
	w.String(&pk.Status)
}

// Unmarshal ...
func (pk *ServerStatusResponse) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Status)
}
