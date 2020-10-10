package event_test

import (
	"testing"

	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/event"
	"github.com/stretchr/testify/assert"
)

func TestTypedApplierCreation(t *testing.T) {
	as := assert.New(t)

	ta := event.TypedApplier{
		"test": func(e *timebox.Event) {},
	}
	applier := ta.Applier()
	as.NotNil(applier)
}

func TestTypedApplierCombine(t *testing.T) {
	as := assert.New(t)
	var results [3]bool
	ta1 := event.TypedApplier{
		"test1": func(_ *timebox.Event) {
			results[0] = true
		},
		"test2": func(_ *timebox.Event) {
			results[1] = true
		},
	}
	ta2 := event.TypedApplier{
		"test1": func(_ *timebox.Event) {
			results[2] = true
		},
	}
	applier := ta1.Combine(ta2).Applier()
	applier(&timebox.Event{
		Type: "test1",
	})
	as.True(results[0])
	as.False(results[1])
	as.True(results[2])

	applier(&timebox.Event{
		Type: "test2",
	})
	as.True(results[0])
	as.True(results[1])
	as.True(results[2])
}
