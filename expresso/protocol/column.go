package protocol

import "fmt"

// ColumnPosition represents a position of a column.
type ColumnPosition [2]int32

// Column represents a chunk column, which contains chunk data, the chunk position, biomes,
// and other useful information for the client.
type Column struct {
	// Position is the position of the column.
	Position ColumnPosition
	// Chunks contain all chunks associated with the column.
	Chunks map[int32]*Chunk
	// Tiles contains all tile entities associated with the column.
	Tiles []map[string]interface{}
	// HeightMaps contains all height maps associated with the column.
	HeightMaps map[string]interface{}
	// Biomes contains all biomes associated with the column.
	Biomes []int32
}

// NewColumn initializes a new empty chunk column.
func NewColumn(pos ColumnPosition) *Column {
	defaultBiomes := make([]int32, 1024)
	for i := 0; i < 1024; i++ {
		defaultBiomes[i] = 1
	}

	return &Column{
		Position:   pos,
		Chunks:     make(map[int32]*Chunk),
		Tiles:      make([]map[string]interface{}, 0),
		HeightMaps: make(map[string]interface{}),
		Biomes:     defaultBiomes,
	}
}

// Get returns the state ID of a block position.
func (c *Column) Get(x, y, z int32) (int32, error) {
	chunk := c.Chunks[y>>4]
	if chunk == nil || chunk.Empty() {
		// The chunk is empty or does not exist, so the block is air.
		return 0, nil
	}

	// Return the state ID from the chunk function.
	return chunk.Get(x, y&15, z)
}

// Set sets the state ID of a block position.
func (c *Column) Set(x, y, z int32, state int32) error {
	chunkIndex := y >> 4
	if chunkIndex < 0 || chunkIndex >= 16 {
		return fmt.Errorf("invalid chunk index")
	}

	chunk, ok := c.Chunks[chunkIndex]
	if !ok {
		if state == air {
			// The chunk does not exist, and we are trying to set a block in the chunk to air,
			// so there is nothing to do.
			return nil
		}

		// Initialize a new empty chunk and update the chunk slice.
		chunk = NewEmptyChunk()
		c.Chunks[chunkIndex] = chunk
	}

	return chunk.Set(x, y&15, z, state)
}
