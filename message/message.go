package message

import (
	"github.com/kode4food/timebox/id"
	"github.com/kode4food/timebox/time"
)

type (
	// ID identifies a Message
	ID = id.ID

	// Message represents the basic form of an Event or Command
	Message struct {
		ID        ID        `json:"id"`
		CreatedAt time.Time `json:"created-at"`
		Type      Type      `json:"type"`
		Payload   Payload   `json:"payload"`
	}

	// Type identifies the message type
	Type string

	// Payload identifies the payload for a message
	Payload = interface{}

	// Handler is an interface that handles an emitted Message
	Handler func(*Message) error
)

// NewID creates a new Message ID
var NewID = id.New

// New returns a new Message given the provided type and payload
func New(msgType Type, payload Payload) *Message {
	return &Message{
		ID:        NewID(),
		CreatedAt: time.Now(),
		Type:      msgType,
		Payload:   payload,
	}
}
