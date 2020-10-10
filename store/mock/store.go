package mock

import (
	"sync"

	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
	"github.com/kode4food/timebox/store/internal/option"
)

type (
	mock struct {
		sync.Mutex
		streams streamMap
		sinkers []store.Sinker
	}

	streamMap map[store.ID]message.List
)

// Open return a new mock Store instance
func Open(o ...option.Option) (store.Store, error) {
	res := &mock{
		streams: streamMap{},
	}
	if err := option.Apply(res, o...); err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mock) SinkTo(s store.Sinker) {
	m.sinkers = append(m.sinkers, s)
}

func (m *mock) sinkEvents(id store.ID, events message.List) {
	for _, s := range m.sinkers {
		s(id, events)
	}
}

func (m *mock) getEvents(id store.ID) message.List {
	if a, ok := m.streams[id]; ok {
		return a
	}
	events := message.EmptyList
	m.streams[id] = events
	return events
}

func (m *mock) putEvents(id store.ID, events message.List) {
	m.streams[id] = events
}

func (m *mock) makeInitialResult(id store.ID) store.Result {
	return &result{
		store: m,
		id:    id,
		next:  store.InitialVersion,
	}
}

func (m *mock) New() (store.Result, error) {
	m.Lock()
	defer m.Unlock()

	id := store.NewID()
	m.streams[id] = message.EmptyList
	return m.makeInitialResult(id), nil
}

func (m *mock) None(id store.ID) (store.Result, error) {
	return m.makeInitialResult(id), nil
}

func (m *mock) All(id store.ID) (store.Result, error) {
	m.Lock()
	defer m.Unlock()

	events := m.getEvents(id)
	return &result{
		store: m,
		id:    id,
		next:  store.Version(len(events)),
	}, nil
}

func (m *mock) Before(id store.ID, v store.Version) (store.Result, error) {
	m.Lock()
	defer m.Unlock()

	events := m.getEvents(id)
	next := store.Version(len(events))
	if v > next {
		return nil, store.NewVersionError(id, next, v)
	}
	return &result{
		store: m,
		id:    id,
		first: store.InitialVersion,
		next:  v,
	}, nil
}
