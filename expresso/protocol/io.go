package protocol

import (
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/text"
)

// IO is implemented by Writer and Reader.
type IO interface {
	// Uint8 reads/writes an uint8 from/to the underlying buffer.
	Uint8(x *uint8)
	// Int16 reads/writes an int16 from/to the underlying buffer.
	Int16(x *int16)
	// Int32 reads/writes an int32 from/to the underlying buffer.
	Int32(x *int32)
	// Int64 reads/writes an int64 from/to the underlying buffer.
	Int64(x *int64)

	// Float32 reads/writes a float32 from/to the underlying buffer.
	Float32(x *float32)
	// Float64 reads/writes a float64 from/to the underlying buffer.
	Float64(x *float64)

	// Varint32 reads/writes a variable int32 from/to the underlying buffer.
	Varint32(x *int32)
	// Varint64 reads/writes a variable int64 from/to the underlying buffer.
	Varint64(x *int64)

	// ByteSlice reads/writes a byte slice from the underlying buffer, similarly to String.
	ByteSlice(x *[]byte)
	// Bytes reads all remaining bytes in the reader, or appends the bytes to the buffer if it is a writer.
	Bytes(b *[]byte)
	// Bool reads/writes a bool as either 0x00 or 0x01 to the underlying buffer.
	Bool(x *bool)
	// String reads/writes a string, prefixed with a variable int32, from/to the underlying buffer.
	String(x *string)

	// UUID reads/writes a UUID from/to the underlying buffer.
	UUID(x *uuid.UUID)
	// Text reads/writes Minecraft-style text from/to the underlying buffer.
	Text(x *text.Text)

	// NBT reads/writes a map as a compound tag from/to the underlying buffer.
	NBT(x *map[string]interface{})
}

// Compile time checks to make sure IO is implemented by Writer and Reader.
var _, _ = IO(&Writer{}), IO(&Reader{})
