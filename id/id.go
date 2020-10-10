package id

import (
	"encoding/json"

	"github.com/google/uuid"
)

// ID is a universally unique identifier
type ID uuid.UUID

// Nil is a the nil ID (all zeroes)
var Nil = ID(uuid.Nil)

// New returns a new universally unique identifier
func New() ID {
	return ID(uuid.New())
}

// Parse parses a string into a valid id
func Parse(idStr string) (ID, error) {
	id, err := uuid.Parse(idStr)
	return ID(id), err
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// MarshalJSON converts an ID to a JSON string
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

// UnmarshalJSON converts a JSON string to an ID
func (id *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	var err error
	*id, err = Parse(s)
	return err
}
