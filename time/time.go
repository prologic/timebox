package time

import (
	"errors"
	"time"
)

// RFC3339NanoSortable is a sortable version of the RFC3339Nano format
const RFC3339NanoSortable = "2006-01-02T15:04:05.000000000Z07:00"

const errRange = "Time.MarshalJSON: year outside of range [0,9999]"

// Time is a value that can represent wall-time
type Time time.Time

// Now returns the current Time using coordinated universal time
func Now() Time {
	return Time(time.Now().UTC())
}

// MarshalJSON converts Time to a JSON string
func (t Time) MarshalJSON() ([]byte, error) {
	ot := time.Time(t)
	if y := ot.Year(); y < 0 || y >= 10000 {
		return nil, errors.New(errRange)
	}

	b := make([]byte, 0, len(RFC3339NanoSortable)+2)
	b = append(b, '"')
	b = ot.AppendFormat(b, RFC3339NanoSortable)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON converts a JSON string to Time
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	ot, err := time.Parse(`"`+RFC3339NanoSortable+`"`, string(data))
	if err != nil {
		return err
	}
	*t = Time(ot)
	return nil
}

func (t Time) String() string {
	ot := time.Time(t)
	b := make([]byte, 0, len(RFC3339NanoSortable)+2)
	b = ot.AppendFormat(b, RFC3339NanoSortable)
	return string(b)
}
