package message

import (
	"encoding/json"

	"github.com/kode4food/timebox/time"
)

// Error messages
const (
	ErrNoTypeDecoder = "decoder not available for type: %s"
)

type (
	// Decoder takes a raw byte stream and decodes it into a Message
	Decoder func([]byte) (*Message, error)

	decodingMessage struct {
		ID        ID              `json:"id"`
		CreatedAt time.Time       `json:"created-at"`
		Type      Type            `json:"type"`
		Payload   json.RawMessage `json:"payload"`
	}
)

// RawDecoder is a decoder that simply returns a Message with a
// json.RawMessage Payload
func RawDecoder(data []byte) (*Message, error) {
	raw := new(decodingMessage)
	if err := json.Unmarshal(data, raw); err != nil {
		return nil, err
	}
	return &Message{
		ID:        raw.ID,
		CreatedAt: raw.CreatedAt,
		Type:      raw.Type,
		Payload:   raw.Payload,
	}, nil
}
