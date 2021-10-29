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
	// to map it. If all else fails, it will return false as it's second return value.
	StateToID(state int32) (int32, bool)
	// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
	// it's second return value.
	IDToState(id int32) (int32, bool)
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
func (*GlobalPalette) StateToID(state int32) (int32, bool) {
	return state, true
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (*GlobalPalette) IDToState(id int32) (int32, bool) {
	return id, true
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
		maxId: maxId,
		data:  make([]int32, maxId+1),
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
	return p.maxId
}

// StateToID converts the block state to a storage ID. If it is not mapped, then the palette will attempt
// to map it. If all else fails, it will return false as it's second return value.
func (p *ListPalette) StateToID(state int32) (id int32, ok bool) {
	for i := int32(0); i < p.nextId; i++ {
		if p.data[i] == state {
			id, ok = i, true
			break
		}
	}

	if !ok && p.Size() < p.maxId+1 {
		p.nextId++

		id, ok = p.nextId, true
		p.data[id] = state
	}

	return id, ok
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (p *ListPalette) IDToState(id int32) (state int32, ok bool) {
	if id >= 0 && id < p.Size() {
		return p.data[id], true
	} else {
		return 0, false
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
func (p *MapPalette) StateToID(state int32) (id int32, ok bool) {
	id, ok = p.stateToID[state]

	if !ok && p.Size() < p.maxId+1 {
		p.nextId++

		id, ok = p.nextId, true
		p.idToState[id] = state
		p.stateToID[state] = id
	}

	return id, ok
}

// IDToState converts the storage ID to a block state. If it is not mapped, then it will return false as
// it's second return value.
func (p *MapPalette) IDToState(id int32) (state int32, ok bool) {
	if id >= 0 && id < p.Size() {
		return p.idToState[id], true
	} else {
		return 0, false
	}
}
