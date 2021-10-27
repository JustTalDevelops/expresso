package protocol

import (
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/text"
)

// IO is implemented by Writer and Reader.
type IO interface {
	// Uint8 reads/writes an uint8 to the underlying buffer.
	Uint8(x *uint8)
	// Int16 writes an int16 to the underlying buffer.
	Int16(x *int16)
	// Int32 writes an int32 to the underlying buffer.
	Int32(x *int32)
	// Int64 writes an Int64 to the underlying buffer.
	Int64(x *int64)

	// Varint32 writes a variable int32 to the underlying buffer.
	Varint32(x *int32)
	// Varint64 writes a variable int64 to the underlying buffer.
	Varint64(x *int64)

	// ByteSlice reads/writes a byte slice from the underlying buffer, similarly to String.
	ByteSlice(x *[]byte)
	// Bytes reads all remaining bytes in the reader/writer.
	Bytes(b *[]byte)
	// Bool reads/writes a bool as either 0x00 or 0x01 to the underlying buffer.
	Bool(x *bool)
	// String writes a string, prefixed with a variable int32, to the underlying buffer.
	String(x *string)
	// UUID reads/writes a UUID from/to the underlying buffer.
	UUID(x *uuid.UUID)
	// Text reads/writes Minecraft-style text from/to the underlying buffer.
	Text(x *text.Text)
}

// Compile time checks to make sure IO is implemented by Writer and Reader.
var _, _ = IO(&Writer{}), IO(&Reader{})
