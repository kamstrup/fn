package seq_test

import (
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestMapSeq(t *testing.T) {
	var arr seq.Seq[int] = seq.MapOf(seq.SliceOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})
	fntesting.TestOf(t, arr).Is(2, 4, 6)
}

func TestMapSeqTake(t *testing.T) {
	head, tail := seq.MapOf[int, int](seq.SliceOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Take(2)

	fntesting.TestOf(t, head.Seq()).Is(2, 4)
	fntesting.TestOf(t, tail).Is(6)
}

func TestMapSeqTakeWhile(t *testing.T) {
	head, tail := seq.MapOf[int, int](seq.SliceOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).TakeWhile(func(i int) bool { return i <= 4 })

	fntesting.TestOf(t, head.Seq()).Is(2, 4)
	fntesting.TestOf(t, tail).Is(6)
}

func TestMapSeqSkip(t *testing.T) {
	tail := seq.MapOf[int, int](seq.SliceOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Skip(1)

	fntesting.TestOf(t, tail).Is(4, 6)
}

func TestMapSeqFirst(t *testing.T) {
	arr := seq.MapOf[int, int](seq.SliceOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})

	var first opt.Opt[int]
	first, arr = arr.First()
	fntesting.TestOf(t, arr).Is(4, 6)
	fntesting.OptOf(t, first).Is(2)

	first, arr = arr.First()
	fntesting.TestOf(t, arr).Is(6)
	fntesting.OptOf(t, first).Is(4)

	first, arr = arr.First()
	fntesting.TestOf(t, arr).IsEmpty()
	fntesting.OptOf(t, first).Is(6)

	first, arr = arr.First()
	fntesting.TestOf(t, arr).IsEmpty()
	fntesting.OptOf(t, first).IsEmpty()
}

func TestMapWhereAnyAll(t *testing.T) {
	m := seq.MapOf[int, int](seq.SliceOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})

	if !seq.All(m, isEven) {
		t.Errorf("all even numbers should be even!")
	}

	if !seq.Any(m, isEven) {
		t.Errorf("all even numbers should be even!")
	}

	if seq.All(m, isOdd) {
		t.Errorf("all even numbers should be even!")
	}

	if seq.Any(m, isOdd) {
		t.Errorf("all even numbers should be even!")
	}

	mz := seq.MapOf[int, int](seq.SliceOfArgs(1, 2, 0, 3), func(i int) int {
		return i * 2
	}).Where(func(i int) bool { return i != 4 })

	if !seq.Any(mz, seq.IsZero[int]) {
		t.Errorf("we should find a zero in mz")
	}
	if !seq.Any(mz, seq.IsNonZero[int]) {
		t.Errorf("we should find a non-zero number in mz")
	}
	if seq.All(mz, seq.IsZero[int]) {
		t.Errorf("mz is not all zeroes!")
	}
	if !seq.Any(mz, seq.IsNonZero[int]) {
		t.Errorf("mz has non-zero elements!")
	}
}
