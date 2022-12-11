package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestZip(t *testing.T) {
	x := fn.ArrayOfArgs(1, 2, 3)
	y := fn.ArrayOfArgs("one", "two", "three", "four")

	z := fn.ZipOf(x, y)
	fntesting.TestOf(t, z).Is(
		fn.TupleOf(1, "one"),
		fn.TupleOf(2, "two"),
		fn.TupleOf(3, "three"),
		// note: seqx is too short, so we truncate seqy
	)
}

func TestZipSource(t *testing.T) {
	x := fn.NumbersFrom(1)
	y := fn.ArrayOfArgs("one", "two", "three", "four")

	z := fn.ZipOf(x, y)
	fntesting.TestOf(t, z).Is(
		fn.TupleOf(1, "one"),
		fn.TupleOf(2, "two"),
		fn.TupleOf(3, "three"),
		fn.TupleOf(4, "four"),
	)
}

func TestZipSuite(t *testing.T) {
	createSeq := func() fn.Seq[fn.Tuple[int, string]] {
		x := fn.NumbersFrom(1)
		y := fn.ArrayOfArgs("one", "two", "three", "four")
		return fn.ZipOf(x, y)
	}

	fntesting.SuiteOf(t, createSeq).Is(
		fn.TupleOf(1, "one"),
		fn.TupleOf(2, "two"),
		fn.TupleOf(3, "three"),
		fn.TupleOf(4, "four"),
	)
}
