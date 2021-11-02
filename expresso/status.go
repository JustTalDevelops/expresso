package expresso

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/text"
)

// Players represents the part of the listener that holds a sample of players.
type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Player `json:"sample"`
}

// Version contains the version and protocol number, used for the status.
type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

// Player represents a player connected to an Expresso listener.
type Player struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

// StatusProvider provides the status of the listener when requested.
type StatusProvider interface {
	// Status returns the status of the listener.
	Status() Status
}

// DefaultStatusProvider is the default status of the listener.
type DefaultStatusProvider struct{}

// Status returns the status of the listener.
func (d *DefaultStatusProvider) Status() Status {
	return Status{
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
	}
}

// Status contains status information about the server. It is used for the multiplayer list.
type Status struct {
	Version     Version   `json:"version"`
	Players     Players   `json:"players"`
	Description text.Text `json:"description"`
	Favicon     string    `json:"favicon,omitempty"`
}

// String returns the status as a string.
func (s Status) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic("should never happen")
	}
	return string(b)
}
