package event

import "github.com/kode4food/timebox/store"

type (
	// Transactional is called with an instantiated Aggregate and any of
	// the Events that are current for that Aggregate at calling time.
	// The Transactional implementation is tasked with retrieving Events
	// from the provided Result, as necessary
	Transactional func(*Aggregate, store.Result) error

	// Source wraps a Store for the purpose of performing transactional
	// event sourcing operations
	Source struct {
		store   store.Store
		handler Handler
	}
)

func noOp(_ *Event) error {
	return nil
}

// NewSource returns a new event sourcing interface
func NewSource(s store.Store) *Source {
	return &Source{
		store:   s,
		handler: noOp,
	}
}

// WithNew creates a new Aggregate in the Store and then invokes the
// provided Transactional function with a reference to that Aggregate
// and any existing Events for it (shouldn't be any). After calling the
// Transactional function, any newly raised Events in the Aggregate
// will be persisted in the Store
func (s *Source) WithNew(t Transactional) error {
	res, err := s.store.New()
	if err != nil {
		return err
	}
	return s.withResult(res, t)
}

// With resolves an existing Aggregate from the Store and then invokes
// the provided Transactional function with a reference to that
// Aggregate and any existing Events for it. After calling the
// Transactional function, any newly raised Events in the Aggregate
// will be persisted in the Store
func (s *Source) With(id store.ID, t Transactional) error {
	res, err := s.store.All(id)
	if err != nil {
		return err
	}
	return s.withResult(res, t)
}

func (s *Source) withResult(r store.Result, t Transactional) error {
	a := NewAggregate(r.ID())
	if err := t(a, r); err != nil {
		return err
	}
	return s.flushTo(a, r)
}

func (s *Source) flushTo(a *Aggregate, r store.Result) error {
	return a.FlushTo(func(l List) error {
		if _, err := r.Append(l); err != nil {
			return err
		}
		if _, err := l.HandleWith(s.handler); err != nil {
			return err
		}
		return nil
	})
}
