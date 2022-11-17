package fn

import (
	"testing"
)

func TestMapSeq(t *testing.T) {
	arr := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	})
	SeqTest[int](t, arr).IsExactly(2, 4, 6)
}

func TestMapSeqTake(t *testing.T) {
	head, tail := SeqMap[int, int](ArrayOfArgs(1, 2, 3), func(i int) int {
		return i * 2
	}).Take(2)

	SeqTest[int](t, head).IsExactly(2, 4)
	SeqTest[int](t, tail).IsExactly(6)
}
