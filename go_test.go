package fn_test

import (
	"strconv"
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestGo(t *testing.T) {
	numStrs := fn.Go(fn.ArrayOfArgs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9), 5, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})
	numStrSorted := numStrs.Array().Sort(fn.OrderAsc[string])
	fntesting.TestOf(t, numStrSorted.Seq()).Is("0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
}
