package expresso

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/protocol/encryption"
	"github.com/justtaldevelops/expresso/expresso/protocol/packet"
	"github.com/justtaldevelops/expresso/expresso/text"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

// Connection is a connection on an Expresso listener.
type Connection struct {
	conn     net.Conn
	listener *Listener

	packets chan packet.Packet

	closed atomic.Bool

	lastKeepAlive atomic.Int64

	threshold atomic.Int32

	packetState atomic.Value

	reader *protocol.Reader
	writer *protocol.Writer

	readMu, writeMu sync.Mutex
}

// defaultCompressionThreshold is always 256.
const defaultCompressionThreshold = 256

// newConn initializes a new Expresso connection.
func newConn(listener *Listener, netConn net.Conn) {
	conn := &Connection{
		conn:     netConn,
		listener: listener,

		packets: make(chan packet.Packet),

		reader: protocol.NewReader(netConn),
		writer: protocol.NewWriter(netConn),
	}
	conn.updateState(packet.StateHandshaking())

	go conn.startReading()
}

// Disconnect disconnects the connection for a given reason.
func (c *Connection) Disconnect(reason text.Text) {
	if c.state() == packet.StateLogin() {
		_ = c.WritePacket(&packet.LoginDisconnect{Reason: reason})
	} else if c.state() == packet.StatePlay() {
		_ = c.WritePacket(&packet.Disconnect{Reason: reason})
	}

	c.Close()
}

// Close closes the connection.
func (c *Connection) Close() {
	c.closed.Store(true)
	_ = c.conn.Close()
}

// WritePacket writes a packet to the connection.
func (c *Connection) WritePacket(pk packet.Packet) error {
	if c.closed.Load() {
		return fmt.Errorf("write packet: connection closed")
	}
	if c.state().Packet(packet.DirectionServer(), pk.ID()) == nil {
		return fmt.Errorf("packet does not exist in current state")
	}

	buf := &bytes.Buffer{}

	w := protocol.NewWriter(buf)
	pk.Marshal(w)

	return c.encode(decodedPacket{id: pk.ID(), contents: buf.Bytes()})
}

// ReadPacket reads a packet from the readable packets available.
func (c *Connection) ReadPacket() (packet.Packet, error) {
	pk, ok := <-c.packets
	if !ok {
		return nil, fmt.Errorf("read packet: connection closed")
	}

	return pk, nil
}

// UpdateCompressionThreshold updates the compression threshold for the connection.
func (c *Connection) UpdateCompressionThreshold(threshold int32) error {
	if threshold != c.CompressionThreshold() {
		// New threshold. Make sure that the client is aware.
		err := c.WritePacket(&packet.SetCompression{Threshold: threshold})
		if err != nil {
			return err
		}
	}

	c.threshold.Store(threshold)
	return nil
}

// CompressionThreshold returns the compression threshold for the connection.
func (c *Connection) CompressionThreshold() int32 {
	return c.threshold.Load()
}

// Compression returns true if the connection is compressing packets.
func (c *Connection) Compression() bool {
	return c.CompressionThreshold() > 0
}

// readPacket reads a packet from a connection.
func (c *Connection) readPacket() (packet.Packet, error) {
	// Decode the newest packet from the connection.
	dec, err := c.decode()
	if err != nil {
		return nil, err
	}

	// Unmarshal it into a packet.
	pk := c.state().Packet(packet.DirectionClient(), dec.id)
	if pk == nil {
		// TODO: Log that there was an unhandled packet
		return c.readPacket()
	}

	pk.Unmarshal(protocol.NewReader(bytes.NewReader(dec.contents)))

	if ok, err := c.handlePacket(pk); ok {
		if err != nil {
			c.Disconnect(text.Text{Text: err.Error(), Color: "red"})
			return nil, fmt.Errorf("read packet when connection closed: %w", err)
		}

		return c.readPacket()
	}

	return pk, nil
}

// updateState updates the connection state.
func (c *Connection) updateState(state packet.State) {
	c.packetState.Store(state)
}

// state returns the connection state.
func (c *Connection) state() packet.State {
	return c.packetState.Load().(packet.State)
}

// keepAlive keeps the connection alive by sending keep alive packets every second.
func (c *Connection) keepAlive() {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	c.lastKeepAlive.Store(time.Now().Unix())

	for {
		select {
		case <-t.C:
			if time.Since(time.Unix(c.lastKeepAlive.Load(), 0)).Seconds() > 30 {
				c.Disconnect(text.Text{Text: "You timed out! Please reconnect and ensure you're not having internet issues."})
				return
			}

			err := c.WritePacket(&packet.ServerKeepAlive{PingID: time.Now().UnixNano() / int64(time.Millisecond)})
			if err != nil {
				c.Close()
				return
			}
		}
	}
}

// startReading starts reading packets from the connection.
func (c *Connection) startReading() {
	for {
		pk, err := c.readPacket()
		if err != nil {
			c.Close()
			break
		}
		if c.state() == packet.StatePlay() {
			c.packets <- pk
		}
	}

	close(c.packets)
}

// handlePacket handles a read packet from the connection.
func (c *Connection) handlePacket(pk packet.Packet) (bool, error) {
	switch pk := pk.(type) {
	case *packet.ClientKeepAlive:
		c.lastKeepAlive.Store(time.Now().Unix())
		return true, nil
	case *packet.Handshake:
		return c.handleHandshake(pk)
	}

	return false, nil
}

// handleHandshake handles the initial handshake.
func (c *Connection) handleHandshake(pk *packet.Handshake) (bool, error) {
	switch pk.NextState {
	case 0x01:
		return c.handlePing()
	case 0x02:
		// Make sure we support the protocol version.
		if pk.Protocol > protocol.CurrentProtocol {
			c.Disconnect(text.Text{Text: fmt.Sprintf("Outdated server! I'm still on %v.", protocol.CurrentMinecraftVersion)})
			return true, nil
		} else if pk.Protocol < protocol.CurrentProtocol {
			c.Disconnect(text.Text{Text: fmt.Sprintf("Outdated client! Please use %v.", protocol.CurrentMinecraftVersion)})
			return true, nil
		}

		// Accept the login.
		return c.handleLogin()
	}
	return false, nil
}

// handlePing handles the server list ping sequence.
func (c *Connection) handlePing() (bool, error) {
	c.updateState(packet.StateStatus())

	for i := 0; i < 2; i++ {
		// Decode the newest packet from the connection.
		pk, err := c.readPacket()
		if err != nil {
			return true, err
		}

		// Handle the part of the sequence we are in.
		switch pk := pk.(type) {
		case *packet.ClientStatusRequest:
			if err = c.WritePacket(&packet.ServerStatusResponse{Status: c.listener.Status().String()}); err != nil {
				return true, err
			}
		case *packet.ClientStatusPing:
			if err = c.WritePacket(&packet.ServerStatusPong{Payload: pk.Payload}); err != nil {
				return true, err
			}
		}
	}

	c.Close()
	return true, nil
}

// handleLogin handles a login attempt from a client.
func (c *Connection) handleLogin() (bool, error) {
	c.updateState(packet.StateLogin())

	// Decode the newest packet from the connection.
	pk, err := c.readPacket()
	if err != nil {
		return true, err
	}
	loginStart := pk.(*packet.LoginStart)

	// Send an encryption request.
	encryptionRequest := &packet.EncryptionRequest{
		PublicKey:   c.listener.keyPair.PublicKey,
		VerifyToken: c.listener.verifyToken,
	}
	err = c.WritePacket(encryptionRequest)
	if err != nil {
		return true, err
	}

	// Get the response from the client.
	pk, err = c.readPacket()
	if err != nil {
		return true, err
	}

	// Decode the shared secret and verify token.
	resp := pk.(*packet.EncryptionResponse)
	sharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, c.listener.keyPair, resp.SharedSecret)
	if err != nil {
		return true, err
	}
	verifyToken, err := rsa.DecryptPKCS1v15(rand.Reader, c.listener.keyPair, resp.VerifyToken)
	if err != nil {
		return true, err
	}

	// Ensure they are valid.
	if len(sharedSecret) != 16 {
		return true, fmt.Errorf("expected shared secret size of 16, instead recieved %v", len(sharedSecret))
	}
	if !bytes.Equal(verifyToken, c.listener.verifyToken) {
		return true, fmt.Errorf("verify tokens do not match")
	}

	// Check what type of UUID we should respond with.
	var uuidForResponse uuid.UUID
	if c.listener.authentication {
		// Make sure that the player is authenticated with Mojang.
		authenticated, data := authenticatedWithMojang(loginStart.Username, sharedSecret, encryptionRequest)
		if !authenticated {
			return true, fmt.Errorf("not authenticated with mojang")
		}

		uuidForResponse = uuid.MustParse(data.UUID)
	} else {
		// The Notchian server does this for whatever reason?
		uuidForResponse, _ = uuid.FromBytes([]byte("OfflinePlayer:" + loginStart.Username))
	}

	// Initialize the new symmetric encryptor.
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return true, err
	}
	c.reader.Reader = cipher.StreamReader{
		S: encryption.NewCFB8Decrypt(block, sharedSecret),
		R: c.conn,
	}
	c.writer.Writer = cipher.StreamWriter{
		S: encryption.NewCFB8Encrypt(block, sharedSecret),
		W: c.conn,
	}

	// Set the default compression.
	err = c.UpdateCompressionThreshold(defaultCompressionThreshold)
	if err != nil {
		return true, err
	}

	// Succeed with login!
	err = c.WritePacket(&packet.LoginSuccess{
		UUID:     uuidForResponse,
		Username: loginStart.Username,
	})
	if err != nil {
		return true, err
	}

	// Play packets can now be used, so we can add it to the listener now.
	c.updateState(packet.StatePlay())
	c.listener.incoming <- c

	go c.keepAlive()

	return true, nil
}
