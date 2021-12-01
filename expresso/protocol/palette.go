package protocol

// This has effectively been ported from Geyser's MCProtocolLib. Thanks a ton!
// https://github.com/GeyserMC/MCProtocolLib

import (
	"math"
)

// Palette is a palette implementation for mapping block states to storage IDs.
type Palette interface {
	// Size returns the known number of block states in the palette.
	Size() int32
	// StateToID converts the block state to a storage ID. If it is not mapped, then the palette will attempt
	// to map it. If all else fails, it will return -1.
	StateToID(state int32) int32
	// IDToState converts the storage ID to a block state. If it is not mapped, then it will return -1.
	IDToState(id int32) int32
}

// PaletteType is a palette type implementation.
type PaletteType struct {
	// MinimumBitsPerEntry is the minimum number of bits per entry for the palette.
	MinimumBitsPerEntry int32
	// MaximumBitsPerEntry is the maximum number of bits per entry for the palette.
	MaximumBitsPerEntry int32
	// StorageSize is the number of bits used to store the palette.
	StorageSize int32
}

// BiomePaletteType returns a biome palette type implementation.
func BiomePaletteType() PaletteType {
	return PaletteType{
		MinimumBitsPerEntry: 1,
		MaximumBitsPerEntry: 3,
		StorageSize:         64,
	}
}

// ChunkPaletteType returns a chunk palette type implementation.
func ChunkPaletteType() PaletteType {
	return PaletteType{
		MinimumBitsPerEntry: 4,
		MaximumBitsPerEntry: 8,
		StorageSize:         4096,
	}
}

// GlobalPalette is a global palette that maps one to one.
type GlobalPalette struct{}

// NewGlobalPalette returns a new global palette.
func NewGlobalPalette() *GlobalPalette {
	return &GlobalPalette{}
}

// Size returns the known number of block states in the palette.
func (*GlobalPalette) Size() int32 {
	return math.MaxInt32
}

// StateToID converts the block state to a storage ID. If it is not mapped, then the palette will attempt
// to map it. If all else fails, it will return false as it's second return value.
func (*GlobalPalette) StateToID(state int32) int32 {
	return state
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (*GlobalPalette) IDToState(id int32) int32 {
	return id
}

// ListPalette is a palette backed by a list.
type ListPalette struct {
	// maxId is the maximum ID that can be mapped.
	maxId int32
	// data contains the block state data.
	data []int32
	// nextId is the next ID to be mapped.
	nextId int32
}

// NewListPalette returns a new list palette.
func NewListPalette(bitsPerEntry int32) *ListPalette {
	maxId := int32((1 << bitsPerEntry) - 1)

	return &ListPalette{
		nextId: 1,
		maxId:  maxId,
		data:   make([]int32, maxId+1),
	}
}

// NewListPaletteFromReader returns a new list palette from the given reader.
func NewListPaletteFromReader(bitsPerEntry int32, reader *Reader) *ListPalette {
	palette := NewListPalette(bitsPerEntry)

	var paletteLength int32
	reader.Varint32(&paletteLength)

	for i := int32(0); i < paletteLength; i++ {
		reader.Varint32(&palette.data[i])
	}

	palette.nextId = paletteLength
	return palette
}

// Size returns the known number of block states in the palette.
func (p *ListPalette) Size() int32 {
	return p.nextId
}

// StateToID converts the block state to a storage ID. If it is not mapped, then the palette will attempt
// to map it. If all else fails, it will return false as it's second return value.
func (p *ListPalette) StateToID(state int32) (id int32) {
	id = -1
	for i := int32(0); i < p.nextId; i++ {
		if p.data[i] == state {
			id = i
			break
		}
	}

	if id == -1 && p.Size() < p.maxId+1 {
		id = p.nextId
		p.data[id] = state

		p.nextId++
	}

	return id
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (p *ListPalette) IDToState(id int32) int32 {
	if id >= 0 && id < p.Size() {
		return p.data[id]
	} else {
		return 0
	}
}

// MapPalette is a palette backed by a map.
type MapPalette struct {
	// maxId is the maximum ID that can be mapped.
	maxId int32
	// nextId is the next ID to be mapped.
	nextId int32

	// idToState is a slice of states, with the index being the storage ID.
	idToState []int32
	// stateToID is a map of states to storage IDs.
	stateToID map[int32]int32
}

// NewMapPalette returns a new map palette.
func NewMapPalette(bitsPerEntry int32) *MapPalette {
	maxId := int32((1 << bitsPerEntry) - 1)

	return &MapPalette{
		nextId:    1,
		maxId:     maxId,
		idToState: make([]int32, maxId+1),
		stateToID: make(map[int32]int32),
	}
}

// NewMapPaletteFromReader returns a new map palette from the given reader.
func NewMapPaletteFromReader(bitsPerEntry int32, reader *Reader) *MapPalette {
	palette := NewMapPalette(bitsPerEntry)

	var paletteLength int32
	reader.Varint32(&paletteLength)

	for i := int32(0); i < paletteLength; i++ {
		var state int32
		reader.Varint32(&state)

		palette.idToState[i] = state
		palette.stateToID[state] = i
	}

	palette.nextId = paletteLength
	return palette
}

// Size returns the known number of block states in the palette.
func (p *MapPalette) Size() int32 {
	return p.nextId
}

// StateToID converts the block state to a storage ID. If it is not mapped, then the palette will attempt
// to map it. If all else fails, it will return false as it's second return value.
func (p *MapPalette) StateToID(state int32) int32 {
	id, ok := p.stateToID[state]
	if !ok && p.Size() < p.maxId+1 {
		id, ok = p.nextId, true

		p.nextId++
		p.idToState[id] = state
		p.stateToID[state] = id
	}

	if !ok {
		return -1
	}
	return id
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (p *MapPalette) IDToState(id int32) int32 {
	if id >= 0 && id < p.Size() {
		return p.idToState[id]
	} else {
		return 0
	}
}

// SingletonPalette is a palette backed by a single block state.
type SingletonPalette struct {
	// state is the block state.
	state int32
}

// NewSingletonPalette returns a new singleton palette.
func NewSingletonPalette(state int32) *SingletonPalette {
	return &SingletonPalette{state: state}
}

// NewSingletonPaletteFromReader returns a new singleton palette from the given reader.
func NewSingletonPaletteFromReader(reader *Reader) *SingletonPalette {
	var state int32
	reader.Int32(&state)

	return NewSingletonPalette(state)
}

// Size returns the known number of block states in the palette.
func (p *SingletonPalette) Size() int32 {
	return 1
}

// StateToID converts the block state to a storage ID. If it is not mapped, then the palette will attempt
// to map it. If all else fails, it will return false as it's second return value.
func (p *SingletonPalette) StateToID(state int32) int32 {
	if p.state == state {
		return 0
	}
	return -1
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (p *SingletonPalette) IDToState(id int32) int32 {
	if id == 0 {
		return p.state
	}
	return 0
}
