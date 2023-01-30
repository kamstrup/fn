package seq_test

import (
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestZip(t *testing.T) {
	x := seq.SliceOfArgs(1, 2, 3)
	y := seq.SliceOfArgs("one", "two", "three", "four")

	z := seq.ZipOf(x, y)
	fntesting.TestOf(t, z).Is(
		seq.TupleOf(1, "one"),
		seq.TupleOf(2, "two"),
		seq.TupleOf(3, "three"),
		// note: seqx is too short, so we truncate seqy
	)
}

func TestZipSource(t *testing.T) {
	x := seq.RangeFrom(1)
	y := seq.SliceOfArgs("one", "two", "three", "four")

	z := seq.ZipOf(x, y)
	fntesting.TestOf(t, z).Is(
		seq.TupleOf(1, "one"),
		seq.TupleOf(2, "two"),
		seq.TupleOf(3, "three"),
		seq.TupleOf(4, "four"),
	)
}

func TestZipSuite(t *testing.T) {
	createSeq := func() seq.Seq[seq.Tuple[int, string]] {
		x := seq.RangeFrom(1)
		y := seq.SliceOfArgs("one", "two", "three", "four")
		return seq.ZipOf(x, y)
	}

	fntesting.SuiteOf(t, createSeq).Is(
		seq.TupleOf(1, "one"),
		seq.TupleOf(2, "two"),
		seq.TupleOf(3, "three"),
		seq.TupleOf(4, "four"),
	)
}
