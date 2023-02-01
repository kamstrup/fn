package seq_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestConcat2x3(t *testing.T) {
	fntesting.SuiteOf(t, func() seq.Seq[int] {
		return seq.ConcatOf(seq.SliceOfArgs(1, 2, 3), seq.SliceOfArgs(4, 5, 6))
	}).Is(1, 2, 3, 4, 5, 6)
}

func TestConcatError(t *testing.T) {
	theError := errors.New("the error")
	cc := seq.ConcatOf(seq.ErrorOf[int](theError))

	if err := seq.Error(cc); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}

func TestConcatWithEmpty(t *testing.T) {
	fntesting.SuiteOf(t, func() seq.Seq[int] {
		return seq.ConcatOf(seq.Empty[int](), seq.SingletOf(1), seq.SliceOfArgs(2, 3))
	}).Is(1, 2, 3)
}
