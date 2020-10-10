package command

import (
	"fmt"

	"github.com/kode4food/timebox/store"
)

// Error messages
const (
	ErrRetryLimitExceeded = "retry limit (%d) exceeded for command %s"
)

// Retry wraps a Command Handler in such a way that it will retry the
// Command whenever a store.VersionConsistencyError is raised, up to
// the specified retry limit
func Retry(maxRetries uint8, handler Handler) Handler {
	retryLimit := int(maxRetries)
	return func(c *Command) error {
		for i := 0; i <= retryLimit; i++ {
			if err := handler(c); err == nil {
				return nil
			} else if _, ok := err.(*store.VersionConsistencyError); !ok {
				return err
			}
		}
		return fmt.Errorf(ErrRetryLimitExceeded, retryLimit, c.ID)
	}
}
