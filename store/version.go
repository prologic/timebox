package store

import (
	"fmt"
	"strconv"
)

type (
	// Version represents the ordering of an event within a stream
	Version uint64

	// VersionConsistencyError is raised when an action would otherwise
	// result in a Version inconsistency. Before() and Append() raise
	// this error. The command package explicitly checks for this error
	// in order to perform dispatch retry logic
	VersionConsistencyError struct {
		Stream   ID
		Expected Version
		Actual   Version
	}
)

// ParseVersion parses a string into a valid Version
func ParseVersion(versionStr string) (Version, error) {
	v, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return Version(v), nil
}

// NewVersionError constructs a new VersionConsistencyError instance
func NewVersionError(id ID, expected Version, actual Version) error {
	return &VersionConsistencyError{
		Stream:   id,
		Expected: expected,
		Actual:   actual,
	}
}

func (e *VersionConsistencyError) Error() string {
	return fmt.Sprintf(ErrBadVersion, e.Stream, e.Expected, e.Actual)
}
