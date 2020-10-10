package lock

import (
	"sync"

	"github.com/kode4food/timebox/store"
)

type (
	// Manager manages locking semantics for a Store
	Manager struct {
		sync.Mutex
		locks
	}

	// Action is a function that produces an intermediate Result
	Action func() (store.Result, error)

	// Locker is used to control read-write access to a ID
	Locker struct {
		sync.RWMutex
		manager *Manager
		id      store.ID
		refs    int
	}

	locks map[store.ID]*Locker
)

// NewManager instantiates a new Lock Manager
func NewManager() *Manager {
	return &Manager{
		locks: locks{},
	}
}

// ForWrite performs an Action during a Write-Lock on the specified
// ID id.
func (m *Manager) ForWrite(id store.ID, a Action) (store.Result, error) {
	al := m.Retain(id)
	defer al.Release()

	al.Lock()
	defer al.Unlock()
	return a()
}

// ForRead performs an Action during a Read-Lock on the specified
// ID id.
func (m *Manager) ForRead(id store.ID, a Action) (store.Result, error) {
	al := m.Retain(id)
	defer al.Release()

	al.RLock()
	defer al.RUnlock()
	return a()
}

// Retain a Locker for the specified ID ID. Must be Released!
func (m *Manager) Retain(id store.ID) *Locker {
	m.Lock()
	defer m.Unlock()

	if al, ok := m.locks[id]; ok {
		al.refs = al.refs + 1
		return al
	}
	al := &Locker{
		manager: m,
		id:      id,
		refs:    1,
	}
	m.locks[id] = al
	return al
}

// Release a Locker for the specified ID id
func (m *Manager) Release(id store.ID) {
	m.Lock()
	defer m.Unlock()

	if al, ok := m.locks[id]; ok {
		al.refs = al.refs - 1
		if al.refs <= 0 {
			delete(m.locks, al.id)
		}
	}
}

// Release this Locker
func (l *Locker) Release() {
	l.manager.Release(l.id)
}
