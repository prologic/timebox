package command_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/timebox/command"
	"github.com/kode4food/timebox/message"
	"github.com/stretchr/testify/assert"
)

func TestTypedHandlerNonError(t *testing.T) {
	as := assert.New(t)

	handle := (command.TypedHandler{
		"test": func(_ *command.Command) error {
			return nil
		},
	}).Handler()

	err := handle(message.New("test", nil))
	as.Nil(err)
}

func TestTypedHandlerError(t *testing.T) {
	as := assert.New(t)

	handle := (command.TypedHandler{
		"test": func(_ *command.Command) error {
			return fmt.Errorf("oopsie")
		},
	}).Handler()

	err := handle(message.New("test", nil))
	as.NotNil(err)
	as.Equal(fmt.Errorf("oopsie"), err)
}

func TestTypedHandlerNotFound(t *testing.T) {
	as := assert.New(t)
	handle := (command.TypedHandler{}).Handler()
	err := handle(message.New("test", nil))
	as.NotNil(err)
	as.Equal(fmt.Errorf(command.ErrTypeNotRegistered, "test"), err)
}

func TestTypedHandlerCombine(t *testing.T) {
	as := assert.New(t)
	th1 := command.TypedHandler{
		"test1": func(_ *command.Command) error {
			return fmt.Errorf("test1")
		},
		"test2": func(_ *command.Command) error {
			return nil
		},
	}
	th2 := command.TypedHandler{
		"test1": func(_ *command.Command) error {
			return nil
		},
	}
	handler := th1.Combine(th2).Handler()
	as.Nil(handler(&command.Command{
		Type: "test1",
	}))
	as.Nil(handler(&command.Command{
		Type: "test2",
	}))
	as.NotNil(handler(&command.Command{
		Type: "test3",
	}))
}
