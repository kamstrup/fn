package seq_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
	fntesting "github.com/kamstrup/fn/testing"
)

func TestValuesSuite(t *testing.T) {
	createSeq := func() seq.Seq[int] {
		optInts := seq.SliceOfArgs(opt.Of(1), opt.Of(2), opt.Of(3))
		return seq.ValuesOf(optInts)
	}

	fntesting.SuiteOf(t, createSeq).Is(1, 2, 3)
}

func TestValuesError(t *testing.T) {
	theError := errors.New("the error")
	optInts := seq.SliceOfArgs(opt.Of(1), opt.Of(2), opt.Of(3), opt.ErrorOf[int](theError))
	vals := seq.ValuesOf(optInts)

	tail := vals.ForEach(func(_ int) {})
	if theError != seq.Error(tail) {
		t.Fatalf("expected error tail: %v", tail)
	}

	tail = vals.ForEachIndex(func(idx, n int) {
		if idx+1 != n {
			t.Fatalf("invalid number at index %d: %d", idx, n)
		}
	})
	if theError != seq.Error(tail) {
		t.Fatalf("expected error tail: %v", tail)
	}

	_, tail = vals.Take(17)
	if theError != seq.Error(tail) {
		t.Fatalf("expected error tail: %v", tail)
	}

	tail = vals.Skip(17)
	if theError != seq.Error(tail) {
		t.Fatalf("expected error tail: %v", tail)
	}
}
