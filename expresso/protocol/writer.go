package protocol

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/text"
	"io"
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
	i := uint16(*x)
	_, _ = w.Write([]byte{byte(i >> 8), byte(i)})
}

// Int32 writes an int32 to the underlying buffer.
func (w *Writer) Int32(x *int32) {
	i := uint32(*x)
	_, _ = w.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
}

// Int64 writes an int64 to the underlying buffer.
func (w *Writer) Int64(x *int64) {
	i := *x
	_, _ = w.Write([]byte{
		byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32),
		byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i),
	})
}

// Varint32 writes a variable int32 to the underlying buffer.
func (w *Writer) Varint32(x *int32) {
	varInt := make([]byte, 0, 5)
	num := uint32(*x)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		varInt = append(varInt, byte(b))
		if num == 0 {
			break
		}
	}

	_, _ = w.Write(varInt)
}

// Varint64 writes a variable int64 to the underlying buffer.
func (w *Writer) Varint64(x *int64) {
	varInt := make([]byte, 0, 10)
	num := uint64(*x)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		varInt = append(varInt, byte(b))
		if num == 0 {
			break
		}
	}

	_, _ = w.Write(varInt)
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
