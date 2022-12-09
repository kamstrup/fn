package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestRangeOfInt(t *testing.T) {
	fntesting.TestOf(t, fn.RangeOf(-3, 2)).Is(-3, -2, -1, 0, 1)
	fntesting.TestOf(t, fn.RangeOf(0, 2)).Is(0, 1)
	fntesting.TestOf(t, fn.RangeOf(27, 27)).IsEmpty()

	fntesting.TestOf(t, fn.RangeOf(2, -3)).Is(2, 1, 0, -1, -2)
	fntesting.TestOf(t, fn.RangeOf(2, 0)).Is(2, 1)

	fntesting.TestOf(t, fn.RangeStepOf(0, 10, 2)).Is(0, 2, 4, 6, 8)
	fntesting.TestOf(t, fn.RangeStepOf(1, 10, 2)).Is(1, 3, 5, 7, 9)

	fntesting.TestOf(t, fn.RangeStepOf(10, 0, 2)).Is(10, 8, 6, 4, 2)
	fntesting.TestOf(t, fn.RangeStepOf(9, 0, 2)).Is(9, 7, 5, 3, 1)
}

func TestRangeOfUint(t *testing.T) {
	fntesting.TestOf(t, fn.RangeOf(uint32(1), uint32(3))).Is(1, 2)
	fntesting.TestOf(t, fn.RangeOf(uint32(10), uint32(8))).Is(10, 9)

	fntesting.TestOf(t, fn.RangeStepOf(uint32(10), uint32(0), uint32(2))).Is(10, 8, 6, 4, 2)
	fntesting.TestOf(t, fn.RangeStepOf(uint32(9), uint32(0), uint32(2))).Is(9, 7, 5, 3, 1)
}

func TestRangeOfIntSuite(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[int] {
		return fn.RangeOf(-2, 2)
	}).Is(-2, -1, 0, 1)
}

func TestRangeOfUintSuite(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[uint32] {
		return fn.RangeStepOf(uint32(10), uint32(0), uint32(3))
	}).Is(10, 7, 4, 1)
}
