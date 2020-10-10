package message_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/timebox/message"
	"github.com/stretchr/testify/assert"
)

var testList = message.List{
	message.New("int", 1),
	message.New("int", 2),
	message.New("int", 3),
	message.New("int", 4),
}

func TestNewEvent(t *testing.T) {
	as := assert.New(t)

	ev1 := message.New("string", "1234")
	as.Equal("1234", ev1.Payload)
}

func TestHandleList(t *testing.T) {
	as := assert.New(t)

	ev := testList
	sum := 0
	rest, err := ev.HandleWith(func(m *message.Message) error {
		sum = sum + m.Payload.(int)
		return nil
	})
	as.Nil(err)
	as.Equal(10, sum)
	as.Equal(0, len(rest))
}

func TestHandleListError(t *testing.T) {
	as := assert.New(t)

	ev := testList
	first := true
	rest, err := ev.HandleWith(func(m *message.Message) error {
		if first {
			first = false
			return nil
		}
		return fmt.Errorf("not the first")
	})

	as.EqualError(err, "not the first")
	as.Equal(3, len(rest))
}
