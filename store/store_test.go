package store_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/timebox/store"
	"github.com/stretchr/testify/assert"
)

func TestVersionErrorString(t *testing.T) {
	as := assert.New(t)
	id := store.NewID()
	err := store.NewVersionError(id, 0, 50)
	as.NotNil(err)
	expected := fmt.Sprintf(store.ErrBadVersion, id, 0, 50)
	as.Equal(expected, err.Error())
}
