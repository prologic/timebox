package event_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/store"
	"github.com/stretchr/testify/assert"
)

func newAggregate(applier event.Applier) *event.Aggregate {
	id := store.NewID()
	a := event.NewAggregate(id)
	a.ApplyTo(applier)
	return a
}

func TestAggregateID(t *testing.T) {
	as := assert.New(t)
	a := newAggregate(nil)
	as.NotEmpty(a.ID())
}

func TestApply(t *testing.T) {
	as := assert.New(t)
	var value int
	a := newAggregate(func(e *timebox.Event) {
		value = value + e.Payload.(int)
	})
	a.Raise(event.New("int", 1))
	as.Equal(1, value)
	a.Raise(event.New("int", 1))
	a.Raise(event.New("int", 1))
	as.Equal(3, value)
}

func TestHydrate(t *testing.T) {
	as := assert.New(t)
	var value int
	a := newAggregate(func(e *timebox.Event) {
		value = value + e.Payload.(int)
	})
	a.HydrateFrom(event.List{
		event.New("int", 1),
		event.New("int", 2),
		event.New("int", 3),
	})
	as.Equal(6, value)
}

func TestFlush(t *testing.T) {
	as := assert.New(t)
	var value int
	a := newAggregate(func(e *timebox.Event) {
		value = value + e.Payload.(int)
	})
	a.Raise(event.New("int", 1))
	a.Raise(event.New("int", 1))
	a.Raise(event.New("int", 1))
	err := a.FlushTo(func(l event.List) error {
		as.Equal(3, len(l))
		as.Equal(1, l[0].Payload.(int))
		as.Equal(1, l[1].Payload.(int))
		as.Equal(1, l[2].Payload.(int))
		return nil
	})
	as.Nil(err)

	err = a.FlushTo(func(l event.List) error {
		as.Equal(0, len(l))
		return nil
	})
}

func TestBadFlush(t *testing.T) {
	as := assert.New(t)
	var value int
	a := newAggregate(func(e *timebox.Event) {
		value = value + e.Payload.(int)
	})
	a.Raise(event.New("int", 1))
	a.Raise(event.New("int", 1))
	a.Raise(event.New("int", 1))
	err := a.FlushTo(func(_ event.List) error {
		return fmt.Errorf("uh-oh")
	})
	as.EqualError(err, "uh-oh")

	err = a.FlushTo(func(l event.List) error {
		as.Equal(3, len(l))
		return nil
	})
}
