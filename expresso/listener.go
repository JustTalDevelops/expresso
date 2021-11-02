package expresso

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"go.uber.org/atomic"
	"log"
	"net"
	"os"
)

// ListenConfig configures certain parts of the listener.
type ListenConfig struct {
	// ErrorLog is a log.Logger that errors that occur during packet handling of clients are written to. By
	// default, ErrorLog is set to one equal to the global logger.
	ErrorLog *log.Logger
	// DisableAuthentication is true if logins should not be verified with Minecraft/Mojang.
	DisableAuthentication bool
	// StatusProvider represents the server list status which is displayed on the multiplayer screen.
	StatusProvider StatusProvider
}

// Listener is an Expresso listener. It listens on TCP for Minecraft packets, decodes them, and allows
// other parts of the program to handle packets.
type Listener struct {
	address        string
	authentication bool

	errorLog *log.Logger

	listener net.Listener

	incoming chan *Connection

	status atomic.Value

	keyPair     *rsa.PrivateKey
	verifyToken []byte
}

// Listen listens on the address provided.
func (cfg ListenConfig) Listen(address string) (*Listener, error) {
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

	if cfg.ErrorLog == nil {
		cfg.ErrorLog = log.New(os.Stderr, "", log.LstdFlags)
	}
	if cfg.StatusProvider == nil {
		cfg.StatusProvider = &DefaultStatusProvider{}
	}

	list := &Listener{address: address, authentication: !cfg.DisableAuthentication, errorLog: cfg.ErrorLog, listener: l, keyPair: key, verifyToken: token, incoming: make(chan *Connection)}
	list.status.Store(cfg.StatusProvider)

	go list.startListening()

	return list, nil
}

// Listen listens with a default listener configuration.
func Listen(address string) (*Listener, error) {
	return ListenConfig{}.Listen(address)
}

// Close closes the listener.
func (l *Listener) Close() {
	_ = l.listener.Close()
	close(l.incoming)
}

// Accept accepts a new connection from the listener.
func (l *Listener) Accept() (*Connection, error) {
	conn, ok := <-l.incoming
	if !ok {
		return nil, fmt.Errorf("listener closed")
	}

	return conn, nil
}

// UpdateStatusProvider updates the server status.
func (l *Listener) UpdateStatusProvider(status StatusProvider) {
	l.status.Store(status)
}

// StatusProvider returns the listeners' status provider.
func (l *Listener) StatusProvider() StatusProvider {
	return l.status.Load().(StatusProvider)
}

// Status returns the server status.
func (l *Listener) Status() Status {
	return l.StatusProvider().Status()
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
}
