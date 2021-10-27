package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// ClientStatusRequest is sent by the client to the server to request the server status.
type ClientStatusRequest struct{}

// ID ...
func (*ClientStatusRequest) ID() int32 {
	return 0x00
}

// Marshal ...
func (*ClientStatusRequest) Marshal(*protocol.Writer) {}

// Unmarshal ...
func (*ClientStatusRequest) Unmarshal(*protocol.Reader) {}
