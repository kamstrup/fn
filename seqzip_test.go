package fn

import "testing"

func TestZip(t *testing.T) {
	x := ArrayOfArgs(1, 2, 3).Seq()
	y := ArrayOfArgs("one", "two", "three", "four").Seq()

	z := ZipOf(x, y)
	SeqTest(t, z).Is(
		TupleOf(1, "one"),
		TupleOf(2, "two"),
		TupleOf(3, "three"),
		// note: seqx is too short, so we truncate seqy
	)
}

func TestZipSource(t *testing.T) {
	x := SourceOf(NumbersFrom(1))
	y := ArrayOfArgs("one", "two", "three", "four").Seq()

	z := ZipOf(x, y)
	SeqTest(t, z).Is(
		TupleOf(1, "one"),
		TupleOf(2, "two"),
		TupleOf(3, "three"),
		TupleOf(4, "four"),
	)
}

func TestZipSuite(t *testing.T) {
	createSeq := func() Seq[Tuple[int, string]] {
		x := SourceOf(NumbersFrom(1))
		y := ArrayOfArgs("one", "two", "three", "four").Seq()
		return ZipOf(x, y)
	}

	SeqTestSuite(t, createSeq).Is(
		TupleOf(1, "one"),
		TupleOf(2, "two"),
		TupleOf(3, "three"),
		TupleOf(4, "four"),
	)
}
