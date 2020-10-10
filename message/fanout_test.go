package message_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/timebox/message"
	"github.com/stretchr/testify/assert"
)

func TestFanOutNonError(t *testing.T) {
	as := assert.New(t)

	f := message.FanOutHandler{
		func(m *message.Message) error {
			return nil
		},
	}
	handle := f.Handler()

	err := handle(message.New("test", nil))
	as.Nil(err)
}

func TestFanOutError(t *testing.T) {
	as := assert.New(t)

	f := message.FanOutHandler{
		func(m *message.Message) error {
			return fmt.Errorf("oopsie")
		},
	}
	handle := f.Handler()

	err := handle(message.New("test", nil)).(*message.FanOutError)
	as.NotNil(err)
	as.Equal(fmt.Errorf("oopsie"), err.Errors[0].Error)
}

func TestFanOutMultiError(t *testing.T) {
	as := assert.New(t)
	h1 := func(m *message.Message) error {
		return fmt.Errorf("oopsie 1")
	}
	h2 := func(m *message.Message) error {
		return nil
	}
	h3 := func(m *message.Message) error {
		return fmt.Errorf("oopsie 2")
	}

	f := message.FanOutHandler{h1, h2, h3}
	handle := f.Handler()

	err := handle(message.New("test", nil)).(*message.FanOutError)
	as.NotNil(err)
	as.Equal(2, len(err.Errors))
	as.Equal(fmt.Errorf("oopsie 1"), err.Errors[0].Error)
	as.Equal(fmt.Sprintf("%p", h1), fmt.Sprintf("%p", err.Errors[0].Handler))
	as.Equal(fmt.Errorf("oopsie 2"), err.Errors[1].Error)
	as.Equal(fmt.Sprintf("%p", h3), fmt.Sprintf("%p", err.Errors[1].Handler))
}
