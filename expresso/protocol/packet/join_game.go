package packet

import (
	_ "embed"
	"github.com/justtaldevelops/expresso/expresso/nbt"
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

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

// JoinGame is sent by the server to the client to join a game.
type JoinGame struct {
	// EntityID is the ID of the player joining the game.
	EntityID int32
	// Hardcore is true if hardcore mode is enabled in the game.
	Hardcore bool
	// GameMode is the game mode of the game.
	GameMode byte
	// PreviousGameMode is the player's previous game mode.
	PreviousGameMode byte
	// Worlds contains all worlds on the server.
	Worlds []string
	// World is the name of the world the player is joining.
	World string
	// HashedSeed contains the first eight bytes of the world seed in an SHA-256 hash.
	HashedSeed int64
	// MaxPlayers was once used by the client to draw the player list, but now is ignored.
	MaxPlayers int32
	// ViewDistance is the maximum view distance the client can use. This ranges from two to thirty-two.
	ViewDistance int32
	// ReducedDebugInfo is true if the client should reduce the amount of debug information it shows on the F3 screen.
	ReducedDebugInfo bool
	// EnableRespawnScreen is true if the client should show the respawn screen when the player dies.
	EnableRespawnScreen bool
	// Debug is true if the world is in debug mode.
	Debug bool
	// Flat is true if the world is flat.
	Flat bool
}

// ID ...
func (*JoinGame) ID() int32 {
	return 0x26
}

// Marshal ...
func (pk *JoinGame) Marshal(w *protocol.Writer) {
	w.Int32(&pk.EntityID)
	w.Bool(&pk.Hardcore)
	w.Uint8(&pk.GameMode)
	w.Uint8(&pk.PreviousGameMode)

	worldsLen := int32(len(pk.Worlds))
	w.Varint32(&worldsLen)
	for _, world := range pk.Worlds {
		w.String(&world)
	}

	w.NBT(&dimensionCodec)
	w.NBT(&dimension)
	w.String(&pk.World)
	w.Int64(&pk.HashedSeed)
	w.Varint32(&pk.MaxPlayers)
	w.Varint32(&pk.ViewDistance)
	w.Bool(&pk.ReducedDebugInfo)
	w.Bool(&pk.EnableRespawnScreen)
	w.Bool(&pk.Debug)
	w.Bool(&pk.Flat)
}

// Unmarshal ...
func (pk *JoinGame) Unmarshal(r *protocol.Reader) {
	r.Int32(&pk.EntityID)
	r.Bool(&pk.Hardcore)
	r.Uint8(&pk.GameMode)
	r.Uint8(&pk.PreviousGameMode)

	var worldsLen int32
	r.Varint32(&worldsLen)

	pk.Worlds = make([]string, worldsLen)
	for i := int32(0); i < worldsLen; i++ {
		r.String(&pk.Worlds[i])
	}

	r.NBT(&dimensionCodec)
	r.NBT(&dimension)
	r.String(&pk.World)
	r.Int64(&pk.HashedSeed)
	r.Varint32(&pk.MaxPlayers)
	r.Varint32(&pk.ViewDistance)
	r.Bool(&pk.ReducedDebugInfo)
	r.Bool(&pk.EnableRespawnScreen)
	r.Bool(&pk.Debug)
	r.Bool(&pk.Flat)
}
