package test

import (
	"io"
	"testing"

	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
	"github.com/kode4food/timebox/store/internal/option"
	"github.com/stretchr/testify/assert"
)

var testMessages = message.List{
	message.New("string", "123"),
	message.New("string", "456"),
	message.New("string", "789"),
	message.New("string", "101112"),
}

// MaybeClose a Store if it implements io.Closer
func MaybeClose(t *testing.T, db store.Store) {
	as := assert.New(t)
	if closer, ok := db.(io.Closer); ok {
		as.Nil(closer.Close())
	}
}

// PerformTestStore performs basic exercising of a Store interface
func PerformTestStore(t *testing.T, open store.Opener, o ...option.Option) {
	as := assert.New(t)

	s, err := open(o...)
	as.NotNil(s)
	as.Nil(err)

	id := store.NewID()
	res, err := s.All(id)
	as.Nil(err)
	as.Equal(store.InitialVersion, res.NextVersion())
	res, err = res.Rest()
	as.Nil(err)
	res, err = res.Append(testMessages[0:2])
	as.Nil(err)

	events, err := res.Events()
	testEqualEvents(t, testMessages[0:2], events)

	res, err = res.Append(testMessages[2:4])
	as.Nil(err)

	events, err = res.Events()
	testEqualEvents(t, testMessages[2:], events)

	res, err = s.None(id)
	as.Nil(err)
	as.Equal(store.InitialVersion, res.NextVersion())
	events, err = res.Events()
	as.Equal(0, len(events))

	MaybeClose(t, s)
}

// PerformTestVersionInconsistency makes sure that a Store interface
// properly handles Version consistency checking
func PerformTestVersionInconsistency(
	t *testing.T, open store.Opener, o ...option.Option,
) {
	as := assert.New(t)

	s, err := open(o...)
	as.NotNil(s)
	as.Nil(err)

	a, err := s.All(store.NewID())
	as.Nil(err)
	res, err := a.Append(testMessages[0:3])
	as.Nil(err)

	e3 := message.New("bytes", []byte("789"))
	_, err = res.Append(message.List{e3})
	as.Nil(err)

	_, err = res.Append(message.List{e3})
	as.Equal(store.NewVersionError(a.ID(), 4, 3), err)
	MaybeClose(t, s)
}

// PerformTestPutNothing makes sure that a Store interface properly
// deals with putting of empty Message Lists
func PerformTestPutNothing(
	t *testing.T, open store.Opener, o ...option.Option,
) {
	as := assert.New(t)

	s, err := open(o...)
	as.NotNil(s)
	as.Nil(err)

	a, err := s.All(store.NewID())
	as.NotNil(a)
	as.Nil(err)

	res, err := a.Append(message.EmptyList)
	as.Nil(err)
	ev, err := res.Events()
	as.Nil(err)
	as.Equal(0, len(ev))
	MaybeClose(t, s)
}

// PerformTestBefore uses the Before method to get events up to a
// certain version and then verifies that a call to Rest retrieves
// the remainder
func PerformTestBefore(t *testing.T, open store.Opener, o ...option.Option) {
	as := assert.New(t)

	s, err := open(o...)
	as.NotNil(s)
	as.Nil(err)

	res, _ := s.New()
	res.Append(testMessages)

	id := res.ID()
	res, _ = s.Before(id, store.Version(2))
	events, _ := res.Events()
	as.Equal(store.Version(0), res.FirstVersion())
	as.Equal(store.Version(2), res.NextVersion())
	testEqualEvents(t, testMessages[0:2], events)

	res, _ = res.Rest()
	events, _ = res.Events()
	as.Equal(store.Version(2), res.FirstVersion())
	as.Equal(store.Version(4), res.NextVersion())
	testEqualEvents(t, testMessages[2:], events)

	MaybeClose(t, s)

}

// PerformTestBadBefore makes sure that calls to a Store's Before
// method properly handle Versions that are higher than what is
// currently persisted
func PerformTestBadBefore(t *testing.T, open store.Opener, o ...option.Option) {
	as := assert.New(t)

	s, err := open(o...)
	as.NotNil(s)
	as.Nil(err)

	id := store.NewID()
	a, err := s.Before(id, store.Version(100))
	as.Nil(a)
	as.Equal(store.NewVersionError(id, 0, 100), err)
	MaybeClose(t, s)
}

// PerformSinking makes sure that Sinkers registered via SinkTo
// are invoked at the right times
func PerformSinking(t *testing.T, open store.Opener, o ...option.Option) {
	as := assert.New(t)

	s, err := open(o...)
	as.NotNil(s)
	as.Nil(err)

	var emptyID store.ID
	var checkID store.ID
	s.SinkTo(func(id store.ID, list message.List) {
		if checkID == emptyID {
			checkID = id
		} else {
			as.EqualValues(id, checkID)
		}
		as.Equal(4, len(list))
	})

	s.SinkTo(func(id store.ID, list message.List) {
		if checkID == emptyID {
			checkID = id
		} else {
			as.EqualValues(id, checkID)
		}
		as.Equal(4, len(list))
	})

	r, err := s.New()
	as.NotNil(r)
	as.Nil(err)

	r, _ = r.Append(testMessages)
	as.NotNil(r)
	as.Nil(err)
}

func testEqualEvents(t *testing.T, list1 message.List, list2 message.List) {
	as := assert.New(t)

	as.Equal(len(list1), len(list2))
	for i := 0; i < len(list1); i++ {
		e1, e2 := list1[i], list2[i]
		as.Equal(e1.ID, e2.ID)
		as.Equal(e1.Type, e2.Type)
		as.Equal(e1.CreatedAt, e2.CreatedAt)
	}
}
