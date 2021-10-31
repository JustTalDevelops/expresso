package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/nbt"
	"github.com/justtaldevelops/expresso/expresso/text"
	"io"
	"io/ioutil"
	"math"
)

// Reader is an instance of a protocol reader.
type Reader struct {
	io.Reader
}

// NewReader initializes a new protocol reader using the buffer passed.
func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

// Uint8 reads an uint8 from the underlying buffer.
func (r *Reader) Uint8(x *uint8) {
	b := make([]byte, 1)
	_, _ = r.Read(b)
	*x = b[0]
}

// Int16 reads an int16 from the underlying buffer.
func (r *Reader) Int16(x *int16) {
	b := make([]byte, 2)
	if _, err := r.Read(b); err != nil {
		panic(err)
	}

	*x = int16(b[0])<<8 | int16(b[1])
}

// Int32 reads an int32 from the underlying buffer.
func (r *Reader) Int32(x *int32) {
	b := make([]byte, 4)
	if _, err := r.Read(b); err != nil {
		panic(err)
	}

	*x = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
}

// Int64 reads an int64 from the underlying buffer.
func (r *Reader) Int64(x *int64) {
	b := make([]byte, 8)
	if _, err := r.Read(b); err != nil {
		panic(err)
	}

	*x = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 |
		int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
}

// Float32 reads a float32 from the underlying buffer.
func (r *Reader) Float32(x *float32) {
	var v int32
	r.Int32(&v)

	*x = math.Float32frombits(uint32(v))
}

// Float64 reads a float64 from the underlying buffer.
func (r *Reader) Float64(x *float64) {
	var v int64
	r.Int64(&v)

	*x = math.Float64frombits(uint64(v))
}

// Varint32 reads a variable int32 from the underlying buffer.
func (r *Reader) Varint32(x *int32) {
	var varInt int32
	for size, sec := 0, byte(0x80); sec&0x80 != 0; size++ {
		if size > 5 {
			panic("varint is too big")
		}

		r.Uint8(&sec)
		varInt |= int32(sec&0x7F) << int32(7*size)
	}

	*x = varInt
}

// Varint64 reads a variable int64 from the underlying buffer.
func (r *Reader) Varint64(x *int64) {
	var varInt int64
	for size, sec := 0, byte(0x80); sec&0x80 != 0; size++ {
		if size >= 10 {
			panic("varlong is too big")
		}

		r.Uint8(&sec)
		varInt |= int64(sec&0x7F) << int64(7*size)
	}

	*x = varInt
}

// Bytes reads all remaining bytes in the reader.
func (r *Reader) Bytes(b *[]byte) {
	var err error
	*b, err = ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
}

// ByteSlice reads a byte slice from the underlying buffer, similarly to String.
func (r *Reader) ByteSlice(x *[]byte) {
	var length int32
	r.Varint32(&length)
	l := int(length)
	if l > math.MaxInt32 {
		panic(fmt.Errorf("bytes surpass maximum int32 length"))
	}
	data := make([]byte, l)
	if _, err := r.Read(data); err != nil {
		panic(err)
	}
	*x = data
}

// Bool reads a bool from 0x00 or 0x01 from the underlying buffer.
func (r *Reader) Bool(x *bool) {
	var b byte
	r.Uint8(&b)

	*x = b != 0
}

// String reads a string, prefixed with a variable int32, from the underlying buffer.
func (r *Reader) String(x *string) {
	var length int32
	r.Varint32(&length)

	l := int(length)
	if l > math.MaxInt32 {
		panic(fmt.Errorf("string too long (bad data?)"))
	}

	data := make([]byte, l)
	if _, err := r.Read(data); err != nil {
		panic(err)
	}

	*x = string(data)
}

// UUID reads a UUID from the underlying buffer.
func (r *Reader) UUID(x *uuid.UUID) {
	_, _ = io.ReadFull(r, (*x)[:])
}

// Text reads Minecraft-style text from the underlying buffer.
func (r *Reader) Text(x *text.Text) {
	var s string
	r.String(&s)

	_ = json.Unmarshal([]byte(s), x)
}

// Chunk reads a chunk from the underlying buffer.
func (r *Reader) Chunk(x *Chunk) {
	var blockCount int16
	r.Int16(&blockCount)

	var bitsPerEntry byte
	r.Uint8(&bitsPerEntry)

	palette := readPalette(int32(bitsPerEntry), r)

	var dataSize int32
	r.Varint32(&dataSize)

	data := make([]int64, dataSize)
	for i := int32(0); i < dataSize; i++ {
		r.Int64(&data[i])
	}

	storage, err := NewBitStorageWithData(int32(bitsPerEntry), chunkSize, data)
	if err != nil {
		panic(err)
	}

	*x = Chunk{
		blockCount: int32(blockCount),
		palette:    palette,
		storage:    storage,
	}
}

// NBT reads a map as a compound tag from the underlying buffer.
func (r *Reader) NBT(x *map[string]interface{}) {
	if err := nbt.NewDecoderWithEncoding(r, nbt.BigEndian).Decode(x); err != nil {
		panic(err)
	}
}
