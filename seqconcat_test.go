package fn

import "testing"

func TestConcat2x3(t *testing.T) {
	SeqTestSuite(t, func() Seq[int] {
		return ConcatOfArgs(ArrayOfArgs(1, 2, 3).Seq(), ArrayOfArgs(4, 5, 6).Seq())
	}).Is(1, 2, 3, 4, 5, 6)
}

func TestConcatWithEmpty(t *testing.T) {
	SeqTestSuite(t, func() Seq[int] {
		return ConcatOfArgs(SeqEmpty[int](), SingletOf(1), ArrayOfArgs(2, 3).Seq())
	}).Is(1, 2, 3)
}
