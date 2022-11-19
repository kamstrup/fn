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
