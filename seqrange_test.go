package fn

import "testing"

func TestRangeOfInt(t *testing.T) {
	SeqTest(t, RangeOf(-3, 2)).Is(-3, -2, -1, 0, 1)
	SeqTest(t, RangeOf(0, 2)).Is(0, 1)
	SeqTest(t, RangeOf(27, 27)).IsEmpty()

	SeqTest(t, RangeOf(2, -3)).Is(2, 1, 0, -1, -2)
	SeqTest(t, RangeOf(2, 0)).Is(2, 1)

	SeqTest(t, RangeStepOf(0, 10, 2)).Is(0, 2, 4, 6, 8)
	SeqTest(t, RangeStepOf(1, 10, 2)).Is(1, 3, 5, 7, 9)

	SeqTest(t, RangeStepOf(10, 0, 2)).Is(10, 8, 6, 4, 2)
	SeqTest(t, RangeStepOf(9, 0, 2)).Is(9, 7, 5, 3, 1)
}

func TestRangeOfUint(t *testing.T) {
	SeqTest(t, RangeOf(uint32(1), uint32(3))).Is(1, 2)
	SeqTest(t, RangeOf(uint32(10), uint32(8))).Is(10, 9)

	SeqTest(t, RangeStepOf(uint32(10), uint32(0), uint32(2))).Is(10, 8, 6, 4, 2)
	SeqTest(t, RangeStepOf(uint32(9), uint32(0), uint32(2))).Is(9, 7, 5, 3, 1)
}

func TestRangeOfIntSuite(t *testing.T) {
	SeqTestSuite(t, func() Seq[int] {
		return RangeOf(-2, 2)
	}).Is(-2, -1, 0, 1)
}

func TestRangeOfUintSuite(t *testing.T) {
	SeqTestSuite(t, func() Seq[uint32] {
		return RangeStepOf(uint32(10), uint32(0), uint32(3))
	}).Is(10, 7, 4, 1)
}
