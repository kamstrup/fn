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

	res := vals.ForEach(func(_ int) {})
	if theError != res.Error() {
		t.Fatalf("expected error tail: %v", res)
	}

	res = vals.ForEachIndex(func(idx, n int) {
		if idx+1 != n {
			t.Fatalf("invalid number at index %d: %d", idx, n)
		}
	})
	if theError != res.Error() {
		t.Fatalf("expected error tail: %v", res)
	}

	elems, tail := vals.Take(17)
	if len(elems) != 3 {
		t.Errorf("expected exactly 3 elements, found %v", elems)
	}

	fst, _ := tail.First()
	if theError != fst.Error() {
		t.Fatalf("expected error tail: %v", tail)
	}

	tail = vals.Skip(17)
	fst, _ = tail.First()
	if theError != fst.Error() {
		t.Fatalf("expected error tail: %v", tail)
	}
}
