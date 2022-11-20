package fn

import "testing"

func TestConcat(t *testing.T) {
	SeqTestSuite(t, func() Seq[int] {
		return ConcatOfArgs(ArrayOfArgs(1, 2, 3).Seq(), ArrayOfArgs(4, 5, 6).Seq())
	}).Is(1, 2, 3, 4, 5, 6)
}
