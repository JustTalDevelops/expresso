package protocol

// This has effectively been ported from Geyser's MCProtocolLib. Thanks a ton!
// https://github.com/GeyserMC/MCProtocolLib

import (
	"fmt"
	"math"
)

// magicValues is a list of magic values used for divideMultiply, divideAdd, and divideShift in the BitStorage struct.
var magicValues = []int32{
	-1, -1, 0, math.MinInt32, 0, 0, 1431655765, 1431655765, 0, math.MinInt32,
	0, 1, 858993459, 858993459, 0, 715827882, 715827882, 0, 613566756, 613566756,
	0, math.MinInt32, 0, 2, 477218588, 477218588, 0, 429496729, 429496729, 0,
	390451572, 390451572, 0, 357913941, 357913941, 0, 330382099, 330382099, 0, 306783378,
	306783378, 0, 286331153, 286331153, 0, math.MinInt32, 0, 3, 252645135, 252645135,
	0, 238609294, 238609294, 0, 226050910, 226050910, 0, 214748364, 214748364, 0,
	204522252, 204522252, 0, 195225786, 195225786, 0, 186737708, 186737708, 0, 178956970,
	178956970, 0, 171798691, 171798691, 0, 165191049, 165191049, 0, 159072862, 159072862,
	0, 153391689, 153391689, 0, 148102320, 148102320, 0, 143165576, 143165576, 0,
	138547332, 138547332, 0, math.MinInt32, 0, 4, 130150524, 130150524, 0, 126322567,
	126322567, 0, 122713351, 122713351, 0, 119304647, 119304647, 0, 116080197, 116080197,
	0, 113025455, 113025455, 0, 110127366, 110127366, 0, 107374182, 107374182, 0,
	104755299, 104755299, 0, 102261126, 102261126, 0, 99882960, 99882960, 0, 97612893,
	97612893, 0, 95443717, 95443717, 0, 93368854, 93368854, 0, 91382282, 91382282,
	0, 89478485, 89478485, 0, 87652393, 87652393, 0, 85899345, 85899345, 0,
	84215045, 84215045, 0, 82595524, 82595524, 0, 81037118, 81037118, 0, 79536431,
	79536431, 0, 78090314, 78090314, 0, 76695844, 76695844, 0, 75350303, 75350303,
	0, 74051160, 74051160, 0, 72796055, 72796055, 0, 71582788, 71582788, 0,
	70409299, 70409299, 0, 69273666, 69273666, 0, 68174084, 68174084, 0, math.MinInt32,
	0, 5,
}

// BitStorage implements the compacted data storage used in Chunks used since Minecraft v1.16.
// https://wiki.vg/Chunk_Format
type BitStorage struct {
	// data is the underlying data storage.
	data []int64
	// bitsPerEntry is the number of bits used to store each entry.
	bitsPerEntry int32
	// size is the number of entries in the storage.
	size int32
	// maxValue is the maximum value that can be stored in the storage.
	maxValue int64
	// valuesPerEntry is the number of values that can be stored in a single entry.
	valuesPerEntry int32
	// divideMultiply could be any value in magicValues.
	divideMultiply int64
	// divideAdd could be any value in magicValues.
	divideAdd int64
	// divideShift could be any value in magicValues.
	divideShift int32
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
		bitsPerEntry:   bitsPerEntry,
		size:           size,
		maxValue:       (int64(1) << bitsPerEntry) - int64(1),
		valuesPerEntry: 64 / bitsPerEntry,
	}

	expectedLength := (size + storage.valuesPerEntry - 1) / storage.valuesPerEntry
	if data == nil {
		storage.data = make([]int64, expectedLength)
	} else {
		dataLength := int32(len(data))
		if dataLength != expectedLength {
			return nil, fmt.Errorf("invalid data length of %v, expected %v", dataLength, expectedLength)
		}

		storage.data = data
	}

	magicIndex := 3 * (storage.valuesPerEntry - 1)
	storage.divideMultiply = int64(magicValues[magicIndex])
	storage.divideAdd = int64(magicValues[magicIndex+1])
	storage.divideShift = magicValues[magicIndex+2]

	return storage, nil
}

// Get returns the value at the given index.
func (b *BitStorage) Get(index int32) (int32, error) {
	if index < 0 || index > b.size-1 {
		return 0, fmt.Errorf("index out of data bounds (%v)", index)
	}

	cellIndex := b.cellIndex(index)
	bitIndex := b.bitIndex(index, cellIndex)
	return int32(b.data[cellIndex] >> bitIndex & b.maxValue), nil
}

// Set sets the value at the given index.
func (b *BitStorage) Set(index int32, value int32) error {
	if index < 0 || index > b.size-1 {
		return fmt.Errorf("index out of data bounds (%v)", index)
	}

	if value < 0 || int64(value) > b.maxValue {
		return fmt.Errorf("value cannot be outside of accepted range (%v)", index)
	}

	cellIndex := b.cellIndex(index)
	bitIndex := b.bitIndex(index, cellIndex)
	b.data[cellIndex] = b.data[cellIndex] & ^(b.maxValue<<bitIndex) | (int64(value)&b.maxValue)<<bitIndex

	return nil
}

// cellIndex uses the data in the BitStorage to get the cell index of the given index.
func (b *BitStorage) cellIndex(index int32) int32 {
	return int32(int64(index)*b.divideMultiply + b.divideAdd>>32>>b.divideShift)
}

// bitIndex uses the data in the BitStorage to get the bit index for the provided index and cellIndex.
func (b *BitStorage) bitIndex(index, cellIndex int32) int32 {
	return (index - cellIndex*b.valuesPerEntry) * b.bitsPerEntry
}
