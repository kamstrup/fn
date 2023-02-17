package seq_test

import (
	"reflect"
	"testing"

	"github.com/kamstrup/fn/seq"
	fntesting "github.com/kamstrup/fn/testing"
)

func TestLimitEmpty(t *testing.T) {
	// Assert optimization that we ger an Empty[T] seq when len is 0
	empty := seq.LimitOf(seq.SliceOfArgs(1, 2, 3), 0)
	if !reflect.DeepEqual(seq.Empty[int](), empty) {
		t.Fatalf("must be of empty type: %v", empty)
	}

	// Must be of empty type when input seq has len 0
	empty = seq.LimitOf(seq.SliceOf([]int{}), 27)
	if !reflect.DeepEqual(seq.Empty[int](), empty) {
		t.Fatalf("must be of empty type: %v", empty)
	}
}

func TestLimitLenRestricted(t *testing.T) {
	l := seq.LimitOf(seq.SliceOfArgs(1, 2, 3), 1)
	sz, ok := l.Len()
	if !ok {
		t.Fatalf("len must be well-defined")
	}
	if sz != 1 {
		t.Fatalf("len must be 1: %v", sz)
	}
}

func TestLimitLenUnrestricted(t *testing.T) {
	// if the seq len is < limit we must return the seq itself again
	arr := seq.SliceOfArgs(1, 2, 3)
	l := seq.LimitOf(arr, 27)
	if !reflect.DeepEqual(arr, l) {
		t.Fatalf("seq with well-defined len < limit must return itself")
	}

}

func TestLimitLenUnlimited(t *testing.T) {
	l := seq.LimitOf(seq.Constant(1), 27)
	sz, ok := l.Len()
	if !ok {
		t.Fatalf("len must be well-defined")
	}
	if sz != 27 {
		t.Fatalf("len must be 27: %v", sz)
	}
}

func TestLimitSuite(t *testing.T) {
	createSeq := func() seq.Seq[int] {
		return seq.LimitOf(seq.SliceOfArgs(1, 2, 3), 4)
	}
	fntesting.SuiteOf(t, createSeq).Is(1, 2, 3)
}
