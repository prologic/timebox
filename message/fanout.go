package message

import "fmt"

// Error messages
const (
	ErrFanOutFailure = "%d handler(s) failed"
)

type (
	// FanOutHandler passes an Event into all of its registered Handlers
	FanOutHandler []Handler

	// FanOutError is raised when any underlying Handler returns an error
	FanOutError struct {
		Errors []FanOutErrorEntry
	}

	// FanOutErrorEntry maps a Handler to the error that it returned
	FanOutErrorEntry struct {
		Handler Handler
		Error   error
	}
)

// Handler dispatches to an underlying registered Handlers
func (f FanOutHandler) Handler() Handler {
	return func(m *Message) error {
		var errs []FanOutErrorEntry
		for _, h := range f {
			if err := h(m); err != nil {
				errs = append(errs, FanOutErrorEntry{
					Handler: h,
					Error:   err,
				})
			}
		}
		if len(errs) > 0 {
			return &FanOutError{
				Errors: errs,
			}
		}
		return nil
	}
}

func (e *FanOutError) Error() string {
	return fmt.Sprintf(ErrFanOutFailure, len(e.Errors))
}
