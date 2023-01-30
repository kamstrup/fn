package seq_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestWhile(t *testing.T) {
	nums := seq.RangeFrom(0).
		While(func(i int) bool { return i < 10 }).
		Values().Seq()
	fntesting.TestOf(t, nums).Is(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
}

func TestWhileSuite(t *testing.T) {
	createSeq := func() seq.Seq[int] {
		return seq.RangeFrom(0).
			While(func(i int) bool { return i < 10 })
	}

	fntesting.SuiteOf(t, createSeq).Is(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
}

func TestWhileError(t *testing.T) {
	theError := errors.New("the error")
	wh := seq.ErrorOf[int](theError).While(seq.IsNonZero[int])

	if err := seq.Error(wh); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
