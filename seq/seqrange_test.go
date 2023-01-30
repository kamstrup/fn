package seq_test

import (
	"math"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestRangeOfInt(t *testing.T) {
	fntesting.TestOf(t, seq.RangeOf(-3, 2)).Is(-3, -2, -1, 0, 1)
	fntesting.TestOf(t, seq.RangeOf(0, 2)).Is(0, 1)
	fntesting.TestOf(t, seq.RangeOf(27, 27)).IsEmpty()

	fntesting.TestOf(t, seq.RangeOf(2, -3)).Is(2, 1, 0, -1, -2)
	fntesting.TestOf(t, seq.RangeOf(2, 0)).Is(2, 1)

	fntesting.TestOf(t, seq.RangeStepOf(0, 10, 2)).Is(0, 2, 4, 6, 8)
	fntesting.TestOf(t, seq.RangeStepOf(1, 10, 2)).Is(1, 3, 5, 7, 9)

	fntesting.TestOf(t, seq.RangeStepOf(10, 0, 2)).Is(10, 8, 6, 4, 2)
	fntesting.TestOf(t, seq.RangeStepOf(9, 0, 2)).Is(9, 7, 5, 3, 1)
}

func TestRangeOfUint(t *testing.T) {
	fntesting.TestOf(t, seq.RangeOf(uint32(1), uint32(3))).Is(1, 2)
	fntesting.TestOf(t, seq.RangeOf(uint32(10), uint32(8))).Is(10, 9)

	fntesting.TestOf(t, seq.RangeStepOf(uint32(10), uint32(0), uint32(2))).Is(10, 8, 6, 4, 2)
	fntesting.TestOf(t, seq.RangeStepOf(uint32(9), uint32(0), uint32(2))).Is(9, 7, 5, 3, 1)
}

func TestRangeOfIntSuite(t *testing.T) {
	fntesting.SuiteOf(t, func() seq.Seq[int] {
		return seq.RangeOf(-2, 2)
	}).Is(-2, -1, 0, 1)
}

func TestRangeOfUintSuite(t *testing.T) {
	fntesting.SuiteOf(t, func() seq.Seq[uint32] {
		return seq.RangeStepOf(uint32(10), uint32(0), uint32(3))
	}).Is(10, 7, 4, 1)
}

func TestRangeFromLen(t *testing.T) {
	sz, ok := seq.RangeFrom(0).Len()
	if !ok || sz != math.MaxInt {
		t.Fatalf("bad len: %d, %v", sz, ok)
	}

	sz, ok = seq.RangeFrom(int8(0)).Len()
	if !ok || sz != math.MaxInt8 {
		t.Fatalf("bad len: %d, %v", sz, ok)
	}

	sz, ok = seq.RangeFrom(uint8(0)).Len()
	if !ok || sz != math.MaxUint8 {
		t.Fatalf("bad len: %d, %v", sz, ok)
	}

	sz, ok = seq.RangeFrom(uint16(0)).Len()
	if !ok || sz != math.MaxUint16 {
		t.Fatalf("bad len: %d, %v", sz, ok)
	}

	sz, ok = seq.RangeFrom(uint(0)).Len()
	if !ok || sz != math.MaxInt {
		t.Fatalf("bad len: %d, %v", sz, ok)
	}
}
