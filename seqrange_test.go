package fn

import "testing"

func TestRangeOf(t *testing.T) {
	SeqTest(t, RangeOf(-3, 2)).Is(-3, -2, -1, 0, 1)
	SeqTest(t, RangeOf(0, 2)).Is(0, 1)
	SeqTest(t, RangeOf(27, 27)).IsEmpty()

	SeqTest(t, RangeOf(2, -3)).Is(2, 1, 0, -1, -2)
	SeqTest(t, RangeOf(2, 0)).Is(2, 1)
}
