package fn

import (
	"testing"
)

func TestMapSeq(t *testing.T) {
	var arr Seq[int] = MapOf(ArrayOfArgs(1, 2, 3).Seq(), func(i int) int {
		return i * 2
	})
	SeqTest(t, arr).Is(2, 4, 6)
}

func TestMapSeqTake(t *testing.T) {
	head, tail := MapOf[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Take(2)

	SeqTest(t, head.Seq()).Is(2, 4)
	SeqTest(t, tail).Is(6)
}

func TestMapSeqTakeWhile(t *testing.T) {
	head, tail := MapOf[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).TakeWhile(func(i int) bool { return i <= 4 })

	SeqTest(t, head.Seq()).Is(2, 4)
	SeqTest(t, tail).Is(6)
}

func TestMapSeqSkip(t *testing.T) {
	tail := MapOf[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Skip(1)

	SeqTest(t, tail).Is(4, 6)
}

func TestMapSeqFirst(t *testing.T) {
	arr := MapOf[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})

	var first Opt[int]
	first, arr = arr.First()
	SeqTest(t, arr).Is(4, 6)
	OptTest(t, first).Is(2)

	first, arr = arr.First()
	SeqTest(t, arr).Is(6)
	OptTest(t, first).Is(4)

	first, arr = arr.First()
	SeqTest(t, arr).IsEmpty()
	OptTest(t, first).Is(6)

	first, arr = arr.First()
	SeqTest(t, arr).IsEmpty()
	OptTest(t, first).IsEmpty()
}

func TestMapWhereAnyAll(t *testing.T) {
	m := MapOf[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})

	if !m.All(isEven) {
		t.Errorf("all even numbers should be even!")
	}

	if !m.Any(isEven) {
		t.Errorf("all even numbers should be even!")
	}

	if m.All(isOdd) {
		t.Errorf("all even numbers should be even!")
	}

	if m.Any(isOdd) {
		t.Errorf("all even numbers should be even!")
	}

	mz := MapOf[int, int](ArrayOfArgs(1, 2, 0, 3), func(i int) int {
		return i * 2
	}).Where(func(i int) bool { return i != 4 })

	if !mz.Any(IsZero[int]) {
		t.Errorf("we should find a zero in mz")
	}
	if !mz.Any(IsNonZero[int]) {
		t.Errorf("we should find a non-zero number in mz")
	}
	if mz.All(IsZero[int]) {
		t.Errorf("mz is not all zeroes!")
	}
	if !mz.Any(IsNonZero[int]) {
		t.Errorf("mz has non-zero elements!")
	}
}
