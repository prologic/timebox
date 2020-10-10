package mock

import (
	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
)

type result struct {
	store *mock
	id    store.ID
	first store.Version
	next  store.Version
}

func (r *result) ID() store.ID {
	return r.id
}

func (r *result) FirstVersion() store.Version {
	return r.first
}

func (r *result) NextVersion() store.Version {
	return r.next
}

func (r *result) Rest() (store.Result, error) {
	s := r.store
	s.Lock()
	defer s.Unlock()

	events := s.getEvents(r.id)
	return &result{
		store: s,
		id:    r.id,
		first: r.next,
		next:  store.Version(len(events)),
	}, nil
}

func (r *result) Append(events message.List) (store.Result, error) {
	r.store.Lock()
	defer r.store.Unlock()

	id := r.id
	prev := r.store.getEvents(id)
	next := store.Version(len(prev))
	if r.next != next {
		return nil, store.NewVersionError(id, next, r.next)
	}
	res := append(prev, events...)
	r.store.putEvents(id, res)
	r.store.sinkEvents(id, events)
	return &result{
		store: r.store,
		id:    id,
		first: r.next,
		next:  store.Version(len(res)),
	}, nil
}

func (r *result) Events() (message.List, error) {
	r.store.Lock()
	defer r.store.Unlock()

	events := r.store.getEvents(r.id)
	return events[int(r.first):int(r.next)], nil
}
