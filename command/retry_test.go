package command_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/command"
	"github.com/kode4food/timebox/store"
	"github.com/stretchr/testify/assert"
)

func TestRetryExceeded(t *testing.T) {
	as := assert.New(t)

	handle := command.Retry(5,
		func(c *timebox.Command) error {
			return store.NewVersionError(c.ID, 0, 1)
		},
	)

	c := command.New("test", nil)
	err := handle(c)
	as.NotNil(err)
	expected := fmt.Errorf(command.ErrRetryLimitExceeded, 5, c.ID)
	as.Equal(expected, err)
}

func TestRetrySuccess(t *testing.T) {
	as := assert.New(t)

	var raised bool
	handle := command.Retry(5,
		func(c *timebox.Command) error {
			if !raised {
				raised = true
				return store.NewVersionError(c.ID, 0, 1)
			}
			return nil
		},
	)

	c := command.New("test", nil)
	err := handle(c)
	as.Nil(err)
}

func TestNormalError(t *testing.T) {
	as := assert.New(t)

	handle := command.Retry(5,
		func(_ *timebox.Command) error {
			return fmt.Errorf("oopsie")
		},
	)

	c := command.New("test", nil)
	err := handle(c)
	as.EqualError(err, "oopsie")
}

func TestZeroRetry(t *testing.T) {
	as := assert.New(t)

	var count = 0
	handle := command.Retry(0,
		func(c *timebox.Command) error {
			count++
			return store.NewVersionError(c.ID, 0, 1)
		},
	)

	c := command.New("test", nil)
	err := handle(c)
	as.Equal(1, count)
	as.NotNil(err)
	expected := fmt.Errorf(command.ErrRetryLimitExceeded, 0, c.ID)
	as.Equal(expected, err)
}
