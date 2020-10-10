package local

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
	"github.com/kode4food/timebox/store/internal/lock"
)

type result struct {
	store *local
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

func (r *result) Events() (message.List, error) {
	var res = message.EmptyList
	for i := r.first; i < r.next; i++ {
		key := getStreamKey(r.id, i)
		buf, err := r.store.database.Get(key)
		if err != nil {
			return nil, err
		}
		e, err := r.store.decoder(buf)
		if err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (r *result) Append(events message.List) (store.Result, error) {
	if len(events) == 0 {
		return r, nil
	}
	return r.withWriteLock(func() (store.Result, error) {
		if err := r.checkVersionConsistency(); err != nil {
			return nil, err
		}

		for i, e := range events {
			ver := r.next + store.Version(i)
			key := getStreamKey(r.id, ver)
			if buf, err := json.Marshal(e); err != nil {
				return nil, err
			} else if err := r.store.database.Put(key, buf); err != nil {
				return nil, err
			}
		}

		next := r.next + store.Version(len(events))
		if err := r.setNextVersion(next); err != nil {
			return nil, err
		}
		r.store.sinkEvents(r.id, events)
		return r.makeResultUntil(next), nil
	})
}

func (r *result) Rest() (store.Result, error) {
	return r.withReadLock(func() (store.Result, error) {
		next, err := r.store.getNextVersion(r.id)
		if err != nil {
			return nil, err
		}
		return r.makeResultUntil(next), nil
	})
}

func (r *result) setNextVersion(ver store.Version) error {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, ver); err != nil {
		return err
	}
	key := getVersionKey(r.id)
	if err := r.store.database.Put(key, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (r *result) withWriteLock(action lock.Action) (store.Result, error) {
	return r.store.locks.ForWrite(r.id, action)
}

func (r *result) withReadLock(action lock.Action) (store.Result, error) {
	return r.store.locks.ForRead(r.id, action)
}

func (r *result) checkVersionConsistency() error {
	next, err := r.store.getNextVersion(r.id)
	if err != nil {
		return err
	}
	if r.next != next {
		return store.NewVersionError(r.id, next, r.next)
	}
	return nil
}

func (r *result) makeResultUntil(next store.Version) *result {
	return &result{
		store: r.store,
		id:    r.id,
		first: r.next,
		next:  next,
	}
}
