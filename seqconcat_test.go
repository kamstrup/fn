package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestConcat2x3(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[int] {
		return fn.ConcatOf(fn.ArrayOfArgs(1, 2, 3), fn.ArrayOfArgs(4, 5, 6))
	}).Is(1, 2, 3, 4, 5, 6)
}

func TestConcatWithEmpty(t *testing.T) {
	fntesting.SuiteOf(t, func() fn.Seq[int] {
		return fn.ConcatOf(fn.SeqEmpty[int](), fn.SingletOf(1), fn.ArrayOfArgs(2, 3))
	}).Is(1, 2, 3)
}
