package protocol

// air is the ID of air in Minecraft.
const air = 0

// ChunkSection represents a modern implementation of a Minecraft chunk section.
type ChunkSection struct {
	// blockCount contains the amount of blocks that aren't air in this section.
	blockCount int16
	// chunkData contains the chunk data for this section.
	chunkData *DataPalette
	// biomeData contains the biome data for this section.
	biomeData *DataPalette
}

// NewEmptyChunkSection creates a new empty chunk section.
func NewEmptyChunkSection() *ChunkSection {
	return &ChunkSection{
		chunkData: NewEmptyChunkDataPalette(),
		biomeData: NewBiomeDataPalette(4),
	}
}

// GetBlockState returns the block state for the given block coordinates.
func (c *ChunkSection) GetBlockState(pos BlockPos) (int32, error) {
	return c.chunkData.GetBlockState(pos)
}

// SetBlockState sets the block state for the given block coordinates.
func (c *ChunkSection) SetBlockState(pos BlockPos, state int32) error {
	curr, err := c.chunkData.SetBlockState(pos, state)
	if err != nil {
		return err
	}

	if state != air && curr == air {
		c.blockCount++
	} else if state == air && curr != air {
		c.blockCount--
	}
	return nil
}

// Empty returns true if this chunk section is empty.
func (c *ChunkSection) Empty() bool {
	return c.blockCount == 0
}
