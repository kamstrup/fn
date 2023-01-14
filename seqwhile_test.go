package fn_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestWhile(t *testing.T) {
	nums := fn.NumbersFrom(0).
		While(func(i int) bool { return i < 10 }).
		Array().Seq()
	fntesting.TestOf(t, nums).Is(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
}

func TestWhileSuite(t *testing.T) {
	createSeq := func() fn.Seq[int] {
		return fn.NumbersFrom(0).
			While(func(i int) bool { return i < 10 })
	}

	fntesting.SuiteOf(t, createSeq).Is(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
}

func TestWhileError(t *testing.T) {
	theError := errors.New("the error")
	wh := fn.ErrorOf[int](theError).While(fn.IsNonZero[int])

	if err := fn.Error(wh); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
