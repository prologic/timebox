package store

import (
	"github.com/kode4food/timebox/id"
	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store/internal/option"
)

const (
	// ErrBadVersion describes a VersionConsistencyError
	ErrBadVersion = "bad version for %s - expected %d, got %d"
)

const (
	// InitialVersion is the initial version of any stream
	InitialVersion Version = 0
)

type (
	// ID used by a Store
	ID = id.ID

	// Opener is the function signature for opening a Store
	Opener func(...option.Option) (Store, error)

	// Sinker is a function that is capable of sinking Events that
	// have successfully been persisted
	Sinker func(ID, message.List)

	// Store exposes methods for managing Message aggregations
	Store interface {
		// New allocates a new ID, generating a unique id for it
		// and returning a Result for writing initial Events
		New() (Result, error)

		// SinkTo registers a Sinker for this Store. Multiple Sinkers
		// can be registered per Store. Because a Sinker is invoked
		// whenever a Result is Appended to the Store, it's very
		// important that an implementation not block
		SinkTo(Sinker)

		// None returns a Result that points to absolute none of the
		// previously committed events for the specified ID,
		// such that calling its Rest() method will yield the first
		// events
		None(ID) (Result, error)

		// All returns a Result that points to all currently committed
		// events for the specified ID
		All(ID) (Result, error)

		// Before returns a Result that points to all of the events up
		// to (but excluding) the specified Version. If the specified
		// Version is greater than the underlying ID's next
		// committable Version, a version consistency error will be
		// raised
		Before(ID, Version) (Result, error)
	}

	// Result is an interface for the handling of events from the
	// Store. Its purpose is to manage the internal state of version
	// consistency without having to surface Versions as an explicit
	// API parameter. Upon creation, a Result points at a specific
	// Version of a ID within the Store. All actions performed
	// are relative to that Version. In order to fast-forward, call to
	// Fetch() is required, which returns its own Result.
	Result interface {
		// ID returns the unique identifier of this Result's All
		ID() ID

		// FirstVersion returns the first Version that this Result
		// contains. Important for partial results that are retrieved
		// using the Rest() method
		FirstVersion() Version

		// NextVersion returns the next Version that will attempt to be
		// written relative to this Result
		NextVersion() Version

		// Rest returns a Result of the remaining committed events for
		// this Result's ID. Should a version consistency error
		// occur, one can use the returned Result of this method to
		// perform additional operations, such as an Append()
		Rest() (Result, error)

		// Append attempts to add events at the next ID Version
		// that this Result points to. If additional events have been
		// committed independently, Append will return a version
		// consistency error. If this occurs, one can use Rest() to
		// pull all committed events and then use its Result to
		// attempt another Append()
		Append(message.List) (Result, error)

		// Events returns all events that are referred to by this
		// Result. Keeping in mind that it may only be a partial
		// result
		Events() (message.List, error)
	}
)

// NewID creates a new Store ID
var NewID = id.New
