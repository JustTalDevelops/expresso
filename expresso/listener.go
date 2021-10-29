package expresso

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/text"
	"go.uber.org/atomic"
	"net"
)

// Listener is an Expresso listener. It listens on TCP for Minecraft packets, decodes them, and allows
// other parts of the program to handle packets.
type Listener struct {
	address  string
	listener net.Listener

	incoming chan *Connection

	status atomic.Value

	keyPair     *rsa.PrivateKey
	verifyToken []byte
}

// Listen listens on the address provided.
func Listen(address string) (*Listener, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	token := make([]byte, 4)
	if _, err = rand.Read(token); err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	list := &Listener{address: address, listener: l, keyPair: key, verifyToken: token, incoming: make(chan *Connection)}
	list.status.Store(Status{
		Version: Version{
			Name:     protocol.CurrentVersion,
			Protocol: protocol.CurrentProtocol,
		},
		Players: Players{Online: 9, Max: 10, Sample: []Player{}},
		Description: text.Text{
			Text:   "An Expresso Listener",
			Color:  "gold",
			Bold:   true,
			Italic: true,
		},
	})

	go list.startListening()

	return list, nil
}

// Close closes the listener.
func (l *Listener) Close() {
	_ = l.listener.Close()
}

// Accept accepts a new connection from the listener.
func (l *Listener) Accept() (*Connection, error) {
	conn, ok := <-l.incoming
	if !ok {
		return nil, fmt.Errorf("listener closed")
	}

	return conn, nil
}

// UpdateStatus updates the server status.
func (l *Listener) UpdateStatus(status Status) {
	l.status.Store(status)
}

// Status returns the server status.
func (l *Listener) Status() Status {
	return l.status.Load().(Status)
}

// startListening starts listening on the listener.
func (l *Listener) startListening() {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			break
		}
		newConn(l, conn)
	}

	close(l.incoming)
}
