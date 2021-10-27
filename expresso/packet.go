package expresso

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"io"
)

// decodedPacket contains the id and contents of an encoded packet.
type decodedPacket struct {
	// id is the id of the packet.
	id int32
	// contents contains the contents of the packet.
	contents []byte
}

// encode writes and encodes a packet to the connection from a decodedPacket.
func (c *Connection) encode(pk decodedPacket) error {
	buf := &bytes.Buffer{}
	w := protocol.NewWriter(buf)
	w.Varint32(&pk.id)
	w.Bytes(&pk.contents)

	if c.Compression() {
		rawLen := buf.Len()
		uncompressedLength := int32(rawLen)
		if uncompressedLength > c.CompressionThreshold() {
			// Compress the packet buffer using zlib and our helper function.
			compress(buf)
		} else {
			// If the compression threshold is more than the uncompressed length, then we just send the packet uncompressed.
			uncompressedLength = 0
		}

		newBuf := &bytes.Buffer{}
		newW := protocol.NewWriter(newBuf)
		newW.Varint32(&uncompressedLength)

		totalLength := int32(newBuf.Len() + rawLen)
		c.writer.Varint32(&totalLength)
		c.writer.Varint32(&uncompressedLength)

		if _, err := buf.WriteTo(c.writer); err != nil {
			return err
		}

		return nil
	}

	b := buf.Bytes()
	l := int32(len(b))

	c.writer.Varint32(&l)
	c.writer.Bytes(&b)
	return nil
}

// decode reads and then decodes a packet from the connection into a decodedPacket.
func (c *Connection) decode() (pk decodedPacket, err error) {
	// Read all packet data.
	var length int32
	c.reader.Varint32(&length)
	if length < 1 {
		return decodedPacket{}, fmt.Errorf("packet length too short: %v", length)
	}

	b := make([]byte, length)
	if _, err := io.ReadFull(c.reader, b); err != nil {
		return decodedPacket{}, fmt.Errorf("read content of packet fail: %w", err)
	}
	buf := bytes.NewBuffer(b)
	r := protocol.NewReader(buf)

	// If the packet data is compressed, decompress it.
	if c.Compression() {
		var uncompressedSize int32
		r.Varint32(&uncompressedSize)

		if uncompressedSize > 0 {
			if err = decompress(buf, uncompressedSize); err != nil {
				return decodedPacket{}, err
			}
		}
	}

	// Read the UUID and contents from the reader.
	r.Varint32(&pk.id)
	r.Bytes(&pk.contents)

	return pk, nil
}

// decompress performs decompression on a zlib compressed buffer. The resulting buffer is returned.
func decompress(compressed *bytes.Buffer, uncompressedSize int32) error {
	uncompressedData := make([]byte, uncompressedSize)
	r, err := zlib.NewReader(compressed)
	if err != nil {
		return fmt.Errorf("decompression failure: %v", err)
	}
	defer r.Close()
	_, err = io.ReadFull(r, uncompressedData)
	if err != nil {
		return fmt.Errorf("decompression failure: %v", err)
	}

	*compressed = *bytes.NewBuffer(uncompressedData)

	return nil
}

// compress performs compression on a buffer using zlib. The resulting buffer is returned.
func compress(uncompressed *bytes.Buffer) {
	compressed := bytes.Buffer{}
	w := zlib.NewWriter(&compressed)
	if _, err := uncompressed.WriteTo(w); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	*uncompressed = compressed
}
