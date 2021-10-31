package packet

import (
	_ "embed"
	"github.com/justtaldevelops/expresso/expresso/nbt"
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

// Packet represents a packet that may be sent over a Minecraft network connection.
type Packet interface {
	// ID returns the ID of the packet. All of these identifiers of packets may be found in id.go.
	ID() int32
	// Marshal encodes the packet to its binary representation into buf.
	Marshal(w *protocol.Writer)
	// Unmarshal decodes a serialised packet in buf into the Packet instance. The serialised packet passed
	// into Unmarshal will not have a header in it.
	Unmarshal(r *protocol.Reader)
}

var (
	// dimensionCodec is a compound tag required for the JoinGame packet which currently has an unknown purpose.
	dimensionCodec map[string]interface{}
	// dimension is a compound tag required for the JoinGame packet which defines valid dimensions.
	dimension map[string]interface{}

	//go:embed dimension_codec.nbt
	dimensionCodecData []byte
	//go:embed dimension.nbt
	dimensionData []byte
)

// init initializes the dimensionCodec and dimension maps.
func init() {
	_ = nbt.Unmarshal(dimensionCodecData, &dimensionCodec)
	_ = nbt.Unmarshal(dimensionData, &dimension)
}
