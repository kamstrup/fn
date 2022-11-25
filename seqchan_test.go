package fn

import "testing"

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
	SeqTest(t, ChanOf(ch)).Is(1, 2, 3)

	// ch is now closed
	SeqTest(t, ChanOf(ch)).IsEmpty()
}

func TestChanSuite(t *testing.T) {
	SeqTestSuite(t, func() Seq[int] {
		return ChanOf(chanWithVals(1, 2, 3, 4))
	}).Is(1, 2, 3, 4)
}
