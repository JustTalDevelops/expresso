package protocol

// This has effectively been ported from Geyser's MCProtocolLib. Thanks a ton!
// https://github.com/GeyserMC/MCProtocolLib

import (
	"fmt"
	"math"
)

const (
	// air is the ID of the air block.
	air = 0
	// chunkSize contains all blocks that are in a chunk.
	chunkSize = 4096
	// minimumPaletteBitsPerEntry is the minimum number of bits per entry in the palette.
	minimumPaletteBitsPerEntry = 4
	// maximumPaletteBitsPerEntry is the maximum number of bits per entry in the palette.
	maximumPaletteBitsPerEntry = 8
	// globalPaletteBitsPerEntry is the number of bits per entry in the global palette.
	globalPaletteBitsPerEntry = 14
)

// Column represents a chunk column, which contains chunk data, the chunk position, biomes,
// and other useful information for the client.
type Column struct {
	// X, Z are the column coordinates.
	X, Z int32
	// Chunks contain all chunks associated with the column.
	Chunks []*Chunk
	// Tiles contains all tile entities associated with the column.
	Tiles []map[string]interface{}
	// HeightMaps contains all height maps associated with the column.
	HeightMaps map[string]interface{}
	// Biomes contains all biomes associated with the column.
	Biomes []int32
}

// Chunk is an implementation of the modern Minecraft chunk.
type Chunk struct {
	// blockCount contains the number of blocks in the chunk.
	blockCount int32
	// palette contains the palette of the chunk.
	palette Palette
	// storage contains the bit storage of the chunk.
	storage *BitStorage
}

// NewEmptyChunk creates a new empty chunk.
func NewEmptyChunk() *Chunk {
	return &Chunk{
		palette: NewListPalette(minimumPaletteBitsPerEntry),
		storage: NewEmptyBitStorage(minimumPaletteBitsPerEntry, chunkSize),
	}
}

// Get returns the block state at the given position.
func (c *Chunk) Get(x, y, z int32) (int32, error) {
	id, err := c.storage.Get(index(x, y, z))
	if err != nil {
		return 0, err
	}
	state, ok := c.palette.IDToState(id)
	if !ok {
		return 0, fmt.Errorf("could not find state for id %v", id)
	}

	return state, nil
}

// Set sets the block state at the given position.
func (c *Chunk) Set(x, y, z, state int32) error {
	id, ok := c.palette.StateToID(state)
	if !ok {
		c.resizePalette()

		id, ok = c.palette.StateToID(state)
		if !ok {
			panic("should never happen")
		}
	}

	ind := index(x, y, z)
	curr, err := c.storage.Get(ind)
	if err != nil {
		return err
	}

	if state != air && curr == air {
		c.blockCount++
	} else if state == air && curr != air {
		c.blockCount--
	}

	return c.storage.Set(ind, id)
}

// Empty returns true if the chunk is empty.
func (c *Chunk) Empty() bool {
	return c.blockCount == 0
}

// resizePalette resizes the palette of the chunk.
func (c *Chunk) resizePalette() {
	bitsPerEntry := sanitizeBitsPerEntry(c.storage.bitsPerEntry + 1)
	newPalette := createPalette(bitsPerEntry)
	newStorage := NewEmptyBitStorage(bitsPerEntry, chunkSize)

	for i := int32(0); i < chunkSize; i++ {
		id, _ := c.storage.Get(i)
		state, _ := c.palette.IDToState(id)
		newID, _ := newPalette.StateToID(state)

		_ = newStorage.Set(i, newID)
	}

	c.palette, c.storage = newPalette, newStorage
}

// sanitizeBitsPerEntry sanitizes the bits per entry of the palette.
func sanitizeBitsPerEntry(bitsPerEntry int32) int32 {
	if bitsPerEntry <= maximumPaletteBitsPerEntry {
		return int32(math.Max(minimumPaletteBitsPerEntry, float64(bitsPerEntry)))
	} else {
		return globalPaletteBitsPerEntry
	}
}

// createPalette creates a new palette with the given number of bits per entry.
func createPalette(bitsPerEntry int32) Palette {
	if bitsPerEntry <= maximumPaletteBitsPerEntry {
		return NewListPalette(bitsPerEntry)
	} else if bitsPerEntry <= maximumPaletteBitsPerEntry {
		return NewMapPalette(bitsPerEntry)
	} else {
		return NewGlobalPalette()
	}
}

// readPalette reads the palette from the given reader.
func readPalette(bitsPerEntry int32, reader *Reader) Palette {
	if bitsPerEntry <= minimumPaletteBitsPerEntry {
		return NewListPaletteFromReader(bitsPerEntry, reader)
	} else if bitsPerEntry <= maximumPaletteBitsPerEntry {
		return NewMapPaletteFromReader(bitsPerEntry, reader)
	} else {
		return NewGlobalPalette()
	}
}

// index converts an X, y, and Z to an integer based index.
func index(x, y, z int32) int32 {
	return y<<8 | z<<4 | x
}
