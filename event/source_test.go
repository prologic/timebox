package event_test

import (
	"testing"

	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/store"
	"github.com/kode4food/timebox/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestEventSourcing(t *testing.T) {
	as := assert.New(t)

	db, _ := mock.Open()
	es := event.NewSource(db)

	var value int
	summer := func(e *timebox.Event) {
		value = value + e.Payload.(int)
	}

	var id store.ID

	tx := func(a *event.Aggregate, r store.Result) error {
		if e, err := r.Events(); err != nil {
			return err
		} else {
			a.ApplyTo(summer)
			a.HydrateFrom(e)
			a.Raise(event.New("int", 3))
			a.Raise(event.New("int", 5))
			return nil
		}
	}

	err := es.WithNew(tx)
	as.Nil(err)
	as.Equal(8, value)
	as.NotEmpty(id)

	err = es.With(id, tx)
	as.Nil(err)
	as.Equal(16, value)
}
