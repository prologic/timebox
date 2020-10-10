package local

import (
	"fmt"

	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store/internal/option"
)

// Path returns an Option for local storage paths
func Path(path string) option.Option {
	return func(t option.Target) error {
		if l, ok := t.(*local); ok {
			l.path = path
			return nil
		}
		return fmt.Errorf("option incorrectly applied")
	}
}

// Decoder returns an Option for local Message decoding
func Decoder(decoder message.Decoder) option.Option {
	return func(t option.Target) error {
		if l, ok := t.(*local); ok {
			l.decoder = decoder
			return nil
		}
		return fmt.Errorf("option incorrectly applied")
	}
}
