package event

import (
	"sync"

	"github.com/kode4food/timebox/store"
)

type (
	// Aggregate manages raised Events, applying them to any registered
	// Appliers and eventually flushing them to a Flusher
	Aggregate struct {
		sync.Mutex
		id       store.ID
		appliers []Applier
		raised   List
	}

	// Applier is an interface that accepts Events to be applied
	Applier func(*Event)

	// Flusher is an all-or-nothing interface for moving a Aggregate's
	// newly raised Events into a persistence mechanism of some kind
	Flusher func(List) error
)

// NewAggregate constructs a new Aggregate manager. Does the work of
// accepting raised events, applying new events, and flushing pending
// events
func NewAggregate(id store.ID) *Aggregate {
	return &Aggregate{
		id:       id,
		appliers: []Applier{},
		raised:   EmptyList,
	}
}

// ID returns the Aggregate's unique identifier
func (a *Aggregate) ID() store.ID {
	return a.id
}

// ApplyTo registers a new Applier with this Aggregate
func (a *Aggregate) ApplyTo(applier Applier) {
	a.appliers = append(a.appliers, applier)
}

func (a *Aggregate) apply(event *Event) {
	for _, a := range a.appliers {
		a(event)
	}
}

// Raise adds an event to the newly raised events for a
// Aggregate and then performs an Apply against its Section
// interface
func (a *Aggregate) Raise(event *Event) {
	a.raised = append(a.raised, event)
	a.apply(event)
}

// HydrateFrom hydrates the Aggregate from the provided List
// of Events
func (a *Aggregate) HydrateFrom(events List) {
	for _, event := range events {
		a.apply(event)
	}
}

// FlushTo attempts to flush newly raised Events to the
// provided Flusher. The internal cache of raised Events is
// cleared upon success
func (a *Aggregate) FlushTo(flush Flusher) error {
	if err := flush(a.raised); err != nil {
		return err
	}
	a.raised = EmptyList
	return nil
}
