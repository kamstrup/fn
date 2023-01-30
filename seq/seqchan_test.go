package seq_test

import (
	"testing"

	"github.com/kamstrup/fn/seq"
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
	fntesting.TestOf(t, seq.ChanOf(ch)).Is(1, 2, 3)

	// ch is now closed
	fntesting.TestOf(t, seq.ChanOf(ch)).IsEmpty()
}

func TestChanSuite(t *testing.T) {
	fntesting.SuiteOf(t, func() seq.Seq[int] {
		return seq.ChanOf(chanWithVals(1, 2, 3, 4))
	}).Is(1, 2, 3, 4)
}
