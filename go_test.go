package fn

import (
	"strconv"
	"testing"
)

func TestGo(t *testing.T) {
	numStrs := Go(ArrayOfArgs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Seq(), 5, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})
	numStrSorted := numStrs.Array().Sort(OrderAsc[string])
	SeqTest(t, numStrSorted.Seq()).Is("0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
}