package time_test

import (
	"testing"
	gotime "time"

	"github.com/kode4food/timebox/time"
	"github.com/stretchr/testify/assert"
)

func TestTime_String(t *testing.T) {
	as := assert.New(t)
	now := time.Now()

	gt := gotime.Time(now)
	b := make([]byte, 0, len(time.RFC3339NanoSortable)+2)
	b = gt.AppendFormat(b, time.RFC3339NanoSortable)

	as.Equal(now.String(), string(b))
}
