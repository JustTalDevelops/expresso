package protocol

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/nbt"
	"github.com/justtaldevelops/expresso/expresso/text"
	"io"
	"math"
)

// Writer is an instance of a protocol writer.
type Writer struct {
	io.Writer
}

// NewWriter initializes a new protocol writer using the buffer passed.
func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}

// Uint8 writes an uint8 to the underlying buffer.
func (w *Writer) Uint8(x *uint8) {
	_, _ = w.Write([]byte{*x})
}

// Int16 writes an int16 to the underlying buffer.
func (w *Writer) Int16(x *int16) {
	i := *x
	_, _ = w.Write([]byte{byte((i >> 8) & 0xFF), byte((i >> 0) & 0xFF)})
}

// Int32 writes an int32 to the underlying buffer.
func (w *Writer) Int32(x *int32) {
	i := *x
	_, _ = w.Write([]byte{byte((i >> 24) & 0xFF), byte((i >> 16) & 0xFF), byte((i >> 8) & 0xFF), byte((i >> 0) & 0xFF)})
}

// Int64 writes an int64 to the underlying buffer.
func (w *Writer) Int64(x *int64) {
	i := *x
	_, _ = w.Write([]byte{
		byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32),
		byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i >> 0),
	})
}

// Float32 writes a float32 to the underlying buffer.
func (w *Writer) Float32(x *float32) {
	bits := int32(math.Float32bits(*x))
	w.Int32(&bits)
}

// Float64 writes a float64 to the underlying buffer.
func (w *Writer) Float64(x *float64) {
	bits := int64(math.Float64bits(*x))
	w.Int64(&bits)
}

// Varint32 writes a variable int32 to the underlying buffer.
func (w *Writer) Varint32(x *int32) {
	i := *x
	for (i & ^0x7F) != 0 {
		_, _ = w.Write([]byte{byte((i & 0x7F) | 0x80)})
		i >>= 7
	}

	_, _ = w.Write([]byte{byte(i)})
}

// Varint64 writes a variable int64 to the underlying buffer.
func (w *Writer) Varint64(x *int64) {
	l := *x
	for (l & ^0x7F) != 0 {
		_, _ = w.Write([]byte{byte(int(l&0x7F) | 0x80)})
		l >>= 7
	}

	_, _ = w.Write([]byte{byte(int(l))})
}

// Bytes appends a []byte to the underlying buffer.
func (w *Writer) Bytes(x *[]byte) {
	_, _ = w.Write(*x)
}

// ByteSlice writes a []byte, prefixed with a variable int32, to the underlying buffer.
func (w *Writer) ByteSlice(x *[]byte) {
	l := int32(len(*x))
	w.Varint32(&l)
	_, _ = w.Write(*x)
}

// Bool writes a bool as either 0x00 or 0x01 to the underlying buffer.
func (w *Writer) Bool(x *bool) {
	v := byte(0x00)
	if *x {
		v = 0x01
	}

	_, _ = w.Write([]byte{v})
}

// String writes a string, prefixed with a variable int32, to the underlying buffer.
func (w *Writer) String(x *string) {
	l := int32(len(*x))
	w.Varint32(&l)
	_, _ = w.Write([]byte(*x))
}

// UUID writes a UUID to the underlying buffer.
func (w *Writer) UUID(x *uuid.UUID) {
	_, _ = w.Write(x[:])
}

// Text writes Minecraft-style text to the underlying buffer.
func (w *Writer) Text(x *text.Text) {
	b, _ := json.Marshal(*x)
	s := string(b)
	w.String(&s)
}

// DataPalette writes a data palette to the underlying buffer.
func (w *Writer) DataPalette(x *DataPalette) {
	if _, ok := x.palette.(*SingletonPalette); ok {
		state := x.palette.IDToState(0)
		bitsPerEntry := uint8(0)
		dataLength := int32(0)

		w.Uint8(&bitsPerEntry)
		w.Varint32(&state)
		w.Varint32(&dataLength)
		return
	}

	bitsPerEntry := uint8(x.storage.bitsPerEntry)
	w.Uint8(&bitsPerEntry)

	if _, ok := x.palette.(*GlobalPalette); !ok {
		paletteLength := x.palette.Size()
		w.Varint32(&paletteLength)

		for i := int32(0); i < paletteLength; i++ {
			state := x.palette.IDToState(i)
			w.Varint32(&state)
		}
	}

	data := x.storage.data
	dataLength := int32(len(data))

	w.Varint32(&dataLength)
	for _, v := range data {
		w.Int64(&v)
	}
}

// ChunkSection writes a chunk section to the underlying buffer.
func (w *Writer) ChunkSection(x *ChunkSection) {
	w.Int16(&x.blockCount)
	w.DataPalette(x.chunkData)
	w.DataPalette(x.biomeData)
}

// NBT writes a map as a compound tag to the underlying buffer.
func (w *Writer) NBT(x *map[string]interface{}) {
	if err := nbt.NewEncoderWithEncoding(w, nbt.BigEndian).Encode(*x); err != nil {
		panic(err)
	}
}
