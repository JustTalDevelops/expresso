package protocol

// This has effectively been copied from go-mc. Many thanks for their work.
// https://github.com/Tnze/go-mc

import (
	"fmt"
	"math"
)

// BitStorage implements the compacted data storage format used in chunks since Minecraft v1.16.
// https://wiki.vg/Chunk_Format
type BitStorage struct {
	data []int64
	mask int64

	bitsPerEntry   int32
	valuesPerEntry int32
	size           int32
}

// NewEmptyBitStorage creates a new empty BitStorage.
func NewEmptyBitStorage(bitsPerEntry int32, size int32) *BitStorage {
	storage, err := NewBitStorageWithData(bitsPerEntry, size, nil)
	if err != nil {
		panic("should never happen")
	}
	return storage
}

// NewBitStorageWithData creates a new BitStorage instance with the provided data.
func NewBitStorageWithData(bitsPerEntry int32, size int32, data []int64) (*BitStorage, error) {
	storage := &BitStorage{
		mask:           1<<bitsPerEntry - 1,
		bitsPerEntry:   bitsPerEntry,
		size:           size,
		valuesPerEntry: 64 / bitsPerEntry,
	}
	dataLen := (size + storage.valuesPerEntry - 1) / storage.valuesPerEntry
	if data != nil {
		otherLength := int32(len(data))
		if otherLength != dataLen {
			return nil, fmt.Errorf("invalid data length of %v, expected %v", dataLen, otherLength)
		}

		storage.data = data
	} else {
		storage.data = make([]int64, dataLen)
	}

	return storage, nil
}

// Set sets the value at the given index.
func (b *BitStorage) Set(index, value int32) error {
	if index < 0 || index > b.size-1 || value < 0 || int64(value) > b.mask {
		return fmt.Errorf("index out of data bounds (%v)", index)
	}

	c, offset := b.calculateIndex(index)
	l := b.data[c]

	b.data[c] = l&(b.mask<<offset^math.MaxInt64) | (int64(value)&b.mask)<<offset

	return nil
}

// Get returns the value at the given index.
func (b *BitStorage) Get(index int32) (int32, error) {
	if index < 0 || index > b.size-1 {
		return 0, fmt.Errorf("index out of data bounds (%v)", index)
	}
	c, offset := b.calculateIndex(index)
	l := b.data[c]

	return int32(l >> offset & b.mask), nil
}

// calculateIndex calculates the new index and offset of the given index.
func (b *BitStorage) calculateIndex(index int32) (int32, int32) {
	ind := index / b.valuesPerEntry
	offset := (index - ind*b.valuesPerEntry) * b.bitsPerEntry

	return ind, offset
}
