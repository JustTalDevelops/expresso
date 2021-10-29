package packet

import "github.com/justtaldevelops/expresso/expresso/protocol"

// ServerPlayerPositionRotation is sent by the server to update the player's position.
type ServerPlayerPositionRotation struct {
	// X, Y, Z are the new coordinates of the player.
	X, Y, Z float64
	// Yaw, Pitch are the new rotation of the player.
	Yaw, Pitch float32

	// Flags are the flags sent by the server in a bitfield.
	Flags byte
	// TeleportID is used if the position update was caused by a teleport.
	TeleportID int32
	// DismountVehicle is used if the player should dismount their vehicle.
	DismountVehicle bool
}

// ID ...
func (*ServerPlayerPositionRotation) ID() int32 {
	return 0x38
}

// Marshal ...
func (pk *ServerPlayerPositionRotation) Marshal(w *protocol.Writer) {
	w.Float64(&pk.X)
	w.Float64(&pk.Y)
	w.Float64(&pk.Z)

	w.Float32(&pk.Yaw)
	w.Float32(&pk.Pitch)

	w.Uint8(&pk.Flags)
	w.Varint32(&pk.TeleportID)
	w.Bool(&pk.DismountVehicle)
}

// Unmarshal ...
func (pk *ServerPlayerPositionRotation) Unmarshal(r *protocol.Reader) {
	r.Float64(&pk.X)
	r.Float64(&pk.Y)
	r.Float64(&pk.Z)

	r.Float32(&pk.Yaw)
	r.Float32(&pk.Pitch)

	r.Uint8(&pk.Flags)
	r.Varint32(&pk.TeleportID)
	r.Bool(&pk.DismountVehicle)
}
