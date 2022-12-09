package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
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
	fn.SeqTest(t, fn.ChanOf(ch)).Is(1, 2, 3)

	// ch is now closed
	fn.SeqTest(t, fn.ChanOf(ch)).IsEmpty()
}

func TestChanSuite(t *testing.T) {
	fn.SeqTestSuite(t, func() fn.Seq[int] {
		return fn.ChanOf(chanWithVals(1, 2, 3, 4))
	}).Is(1, 2, 3, 4)
}
