package mock_test

import (
	"testing"

	"github.com/kode4food/timebox/store/mock"
	st "github.com/kode4food/timebox/store/test"
)

func TestStore(t *testing.T) {
	st.PerformTestStore(t, mock.Open)
}

func TestVersionInconsistency(t *testing.T) {
	st.PerformTestVersionInconsistency(t, mock.Open)
}

func TestPutNothing(t *testing.T) {
	st.PerformTestPutNothing(t, mock.Open)
}

func TestBefore(t *testing.T) {
	st.PerformTestBefore(t, mock.Open)
}

func TestBadBefore(t *testing.T) {
	st.PerformTestBadBefore(t, mock.Open)
}

func TestSinking(t *testing.T) {
	st.PerformSinking(t, mock.Open)
}
