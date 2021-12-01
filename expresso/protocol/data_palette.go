package protocol

// This has effectively been ported from Geyser's MCProtocolLib. Thanks a ton!
// https://github.com/GeyserMC/MCProtocolLib

// globalPaletteBitsPerEntry is the number of bitsPerEntry per entry in the global palette.
const globalPaletteBitsPerEntry = 14

// DataPalette is an implementation of the modern Minecraft data palette.
type DataPalette struct {
	// palette contains the palette of the chunk.
	palette Palette
	// storage contains the bit storage of the chunk.
	storage *BitStorage
	// paletteType contains the type of the palette.
	paletteType PaletteType
	// globalPaletteBits contains the size of the storage.
	globalPaletteBits int32
}

// NewEmptyChunkDataPalette creates a new empty chunk data palette.
func NewEmptyChunkDataPalette() *DataPalette {
	return NewChunkDataPalette(globalPaletteBitsPerEntry)
}

// NewChunkDataPalette creates a new chunk data palette with the globalPaletteBits given.
func NewChunkDataPalette(globalPaletteBits int32) *DataPalette {
	return NewEmptyDataPalette(ChunkPaletteType(), globalPaletteBits)
}

// NewBiomeDataPalette creates a new biome data palette with the globalPaletteBits given.
func NewBiomeDataPalette(globalPaletteBits int32) *DataPalette {
	return NewEmptyDataPalette(BiomePaletteType(), globalPaletteBits)
}

// NewEmptyDataPalette creates a new empty chunk.
func NewEmptyDataPalette(paletteType PaletteType, globalPaletteBits int32) *DataPalette {
	return &DataPalette{
		palette:           NewListPalette(paletteType.MinimumBitsPerEntry),
		storage:           NewEmptyBitStorage(paletteType.MinimumBitsPerEntry, paletteType.StorageSize),
		paletteType:       paletteType,
		globalPaletteBits: globalPaletteBits,
	}
}

// GetBlockState returns the block state at the given position.
func (d *DataPalette) GetBlockState(x, y, z int32) (int32, error) {
	if d.storage != nil {
		id, err := d.storage.Get(index(x, y, z))
		if err != nil {
			return 0, err
		}
		return d.palette.IDToState(id), nil
	} else {
		return d.palette.IDToState(0), nil
	}
}

// SetBlockState sets the block state at the given position.
func (d *DataPalette) SetBlockState(x, y, z, state int32) (int32, error) {
	id := d.palette.StateToID(state)
	if id == -1 {
		d.resize()
		id = d.palette.StateToID(state)
	}

	if d.storage != nil {
		ind := index(x, y, z)
		curr, err := d.storage.Get(ind)
		if err != nil {
			return 0, err
		}

		err = d.storage.Set(ind, id)
		if err != nil {
			return 0, err
		}
		return curr, nil
	}

	// Singleton palette and the block has not changed because the palette hasn't resized
	return state, nil
}

// resize performs a resize on the palette of the chunk.
func (d *DataPalette) resize() {
	bitsPerEntry := int32(1)
	if _, ok := d.palette.(*SingletonPalette); !ok {
		bitsPerEntry = d.storage.bitsPerEntry + 1
	}

	bitsPerEntry = d.sanitizeBitsPerEntry(bitsPerEntry)
	newPalette := createPalette(bitsPerEntry, d.paletteType)
	newStorage := NewEmptyBitStorage(bitsPerEntry, d.paletteType.StorageSize)

	if _, ok := d.palette.(*SingletonPalette); ok {
		for i := int32(0); i < d.paletteType.StorageSize; i++ {
			_ = newStorage.Set(i, 0)
		}
	} else {
		for i := int32(0); i < d.paletteType.StorageSize; i++ {
			id, _ := d.storage.Get(i)
			_ = newStorage.Set(i, newPalette.StateToID(d.palette.IDToState(id)))
		}
	}

	d.palette, d.storage = newPalette, newStorage
}

// sanitizeBitsPerEntry sanitizes the bitsPerEntry per entry of the palette.
func (d *DataPalette) sanitizeBitsPerEntry(bitsPerEntry int32) int32 {
	if bitsPerEntry <= d.paletteType.MaximumBitsPerEntry {
		if bitsPerEntry < d.paletteType.MinimumBitsPerEntry {
			return d.paletteType.MinimumBitsPerEntry
		}
		return bitsPerEntry
	} else {
		return globalPaletteBitsPerEntry
	}
}

// readPalette reads the palette from the given reader.
func readPalette(paletteType PaletteType, bitsPerEntry int32, reader *Reader) Palette {
	if bitsPerEntry > paletteType.MaximumBitsPerEntry {
		return NewGlobalPalette()
	}
	if bitsPerEntry == 0 {
		return NewSingletonPaletteFromReader(reader)
	}
	if bitsPerEntry <= paletteType.MinimumBitsPerEntry {
		return NewListPaletteFromReader(bitsPerEntry, reader)
	}
	return NewMapPaletteFromReader(bitsPerEntry, reader)
}

// createPalette creates a new palette with the given number of bitsPerEntry per entry.
func createPalette(bitsPerEntry int32, paletteType PaletteType) Palette {
	if bitsPerEntry <= paletteType.MinimumBitsPerEntry {
		return NewListPalette(bitsPerEntry)
	} else if bitsPerEntry <= paletteType.MaximumBitsPerEntry {
		return NewMapPalette(bitsPerEntry)
	} else {
		return NewGlobalPalette()
	}
}

// index converts an X, y, and Z to an integer based index.
func index(x, y, z int32) int32 {
	return y<<8 | z<<4 | x
}
