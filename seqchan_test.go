package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func chanWithVals[T any](tt ...T) <-chan T {
	ch := make(chan T)
	go func() {
		for _, t := range tt {
			ch <- t
		}
		close(ch)
	}()
	return ch
}

func TestChan(t *testing.T) {
	ch := chanWithVals(1, 2, 3)
	fntesting.TestOf(t, fn.ChanOf(ch)).Is(1, 2, 3)

	// ch is now closed
	fntesting.TestOf(t, fn.ChanOf(ch)).IsEmpty()
}

func TestChanSuite(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[int] {
		return fn.ChanOf(chanWithVals(1, 2, 3, 4))
	}).Is(1, 2, 3, 4)
}
