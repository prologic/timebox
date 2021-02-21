package message

import (
	"encoding/json"
	"fmt"
)

type (
	// Instantiator creates a new value to be decoded
	Instantiator func() Payload

	// TypedInstantiator is how a Decoder creates a Payload based on Type
	TypedInstantiator map[Type]Instantiator
)

func (t TypedInstantiator) Decoder() Decoder {
	return func(data []byte) (*Message, error) {
		raw := new(decodingMessage)
		if err := json.Unmarshal(data, raw); err != nil {
			return nil, err
		}
		instantiator, ok := t[raw.Type]
		if !ok {
			return nil, fmt.Errorf(ErrNoTypeDecoder, raw.Type)
		}
		payload := instantiator()
		if err := json.Unmarshal(raw.Payload, payload); err != nil {
			return nil, err
		}
		return &Message{
			ID:        raw.ID,
			Type:      raw.Type,
			CreatedAt: raw.CreatedAt,
			Payload:   payload,
		}, nil
	}
}

// Combine with other instances, yielding a new instance
func (t TypedInstantiator) Combine(
	instantiators ...TypedInstantiator,
) TypedInstantiator {
	combined := append([]TypedInstantiator{t}, instantiators...)
	return TypedInstantiators(combined...)
}

// TypedInstantiators combines a TypedInstantiator set
func TypedInstantiators(instantiators ...TypedInstantiator) TypedInstantiator {
	res := TypedInstantiator{}
	for _, i := range instantiators {
		for k, v := range i {
			res[k] = v
		}
	}
	return res
}
