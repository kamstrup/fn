package fn_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestConcat2x3(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[int] {
		return fn.ConcatOf(fn.SliceOfArgs(1, 2, 3), fn.SliceOfArgs(4, 5, 6))
	}).Is(1, 2, 3, 4, 5, 6)
}

func TestConcatError(t *testing.T) {
	theError := errors.New("the error")
	cc := fn.ConcatOf(fn.ErrorOf[int](theError))

	if err := fn.Error(cc); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}

func TestConcatWithEmpty(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[int] {
		return fn.ConcatOf(fn.SeqEmpty[int](), fn.SingletOf(1), fn.SliceOfArgs(2, 3))
	}).Is(1, 2, 3)
}
