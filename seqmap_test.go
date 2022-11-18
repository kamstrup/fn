package fn

import (
	"testing"
)

func TestMapSeq(t *testing.T) {
	arr := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})
	SeqTest[int](t, arr).Is(2, 4, 6)
}

func TestMapSeqTake(t *testing.T) {
	head, tail := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Take(2)

	SeqTest[int](t, head).Is(2, 4)
	SeqTest[int](t, tail).Is(6)
}

func TestMapSeqTakeWhile(t *testing.T) {
	head, tail := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).TakeWhile(func(i int) bool { return i <= 4 })

	SeqTest[int](t, head).Is(2, 4)
	SeqTest[int](t, tail).Is(6)
}

func TestMapSeqSkip(t *testing.T) {
	tail := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Skip(1)

	SeqTest[int](t, tail).Is(4, 6)
}

func TestMapSeqFirst(t *testing.T) {
	arr := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})

	var first Opt[int]
	first, arr = arr.First()
	SeqTest[int](t, arr).Is(4, 6)
	OptTest[int](t, first).Is(2)

	first, arr = arr.First()
	SeqTest[int](t, arr).Is(6)
	OptTest[int](t, first).Is(4)

	first, arr = arr.First()
	SeqTest[int](t, arr).IsEmpty()
	OptTest[int](t, first).Is(6)

	first, arr = arr.First()
	SeqTest[int](t, arr).IsEmpty()
	OptTest[int](t, first).IsEmpty()
}
