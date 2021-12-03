package packet

import (
	"bytes"
	"github.com/bits-and-blooms/bitset"
	"github.com/justtaldevelops/expresso/expresso/protocol"
)

// ChunkData is sent by the server to update a chunk client-side
type ChunkData struct {
	// Column is the chunk column that is being referenced.
	Column *protocol.Column
}

// ID ...
func (*ChunkData) ID() int32 {
	return 0x22
}

// Marshal ...
func (pk *ChunkData) Marshal(w *protocol.Writer) {
	// Bit set and chunk writing.
	dataBuffer := &bytes.Buffer{}
	dataWriter := protocol.NewWriter(dataBuffer)

	bitSet := &bitset.BitSet{}

	for index, chunk := range pk.Column.Chunks {
		if !chunk.Empty() {
			bitSet.Set(uint(index))
			dataWriter.Chunk(chunk)
		}
	}

	// DataPalette position.
	w.Int32(&pk.Column.Position[0])
	w.Int32(&pk.Column.Position[1])

	// Write bitset to main writer.
	bitSetIntegers := bitSet.Bytes()
	bitSetSize := int32(len(bitSetIntegers))

	w.Varint32(&bitSetSize)

	for _, setInteger := range bitSetIntegers {
		integer := int64(setInteger)
		w.Int64(&integer)
	}

	// Height maps.
	w.NBT(&pk.Column.HeightMaps)

	// Biomes.
	biomesLength := int32(len(pk.Column.Biomes))
	w.Varint32(&biomesLength)

	for _, biome := range pk.Column.Biomes {
		w.Varint32(&biome)
	}

	// Write data to main writer.
	dataBytes := dataBuffer.Bytes()
	w.ByteSlice(&dataBytes)

	// Tile entities.
	tileEntitiesSize := int32(len(pk.Column.Tiles))
	w.Varint32(&tileEntitiesSize)

	for _, tileEntity := range pk.Column.Tiles {
		w.NBT(&tileEntity)
	}
}

// Unmarshal ...
func (pk *ChunkData) Unmarshal(r *protocol.Reader) {
	// DataPalette position.
	r.Int32(&pk.Column.Position[0])
	r.Int32(&pk.Column.Position[1])

	// Read chunk mask.
	var bitSetSize int32
	r.Varint32(&bitSetSize)

	bits := make([]uint64, bitSetSize)
	for i := 0; i < int(bitSetSize); i++ {
		var bit int64
		r.Int64(&bit)

		bits[i] = uint64(bit)
	}

	chunkMask := bitset.From(bits)

	// Height maps.
	r.NBT(&pk.Column.HeightMaps)

	// Biomes.
	var biomesLength int32
	r.Varint32(&biomesLength)

	pk.Column.Biomes = make([]int32, biomesLength)
	for i := 0; i < int(biomesLength); i++ {
		r.Varint32(&pk.Column.Biomes[i])
	}

	// Data.
	var data []byte
	r.ByteSlice(&data)

	dataReader := protocol.NewReader(bytes.NewReader(data))
	for index := uint(0); index < chunkMask.Count(); index++ {
		if chunkMask.Test(index) {
			var chunk *protocol.DataPalette
			dataReader.Chunk(chunk)

			pk.Column.Chunks[int32(index)] = chunk
		}
	}

	// Tile entities.
	var tileEntitiesSize int32
	r.Varint32(&tileEntitiesSize)

	pk.Column.Tiles = make([]map[string]interface{}, tileEntitiesSize)
	for i := 0; i < int(tileEntitiesSize); i++ {
		r.NBT(&pk.Column.Tiles[i])
	}
}
