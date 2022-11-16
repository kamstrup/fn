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
