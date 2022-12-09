package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
)

func TestSource(t *testing.T) {
	nums, tail := fn.NumbersFrom(0).
		Map(func(i int) int { return i * 2 }).
		Take(3)
	fn.SeqTest(t, nums.Seq()).Is(0, 2, 4)

	nums, tail = tail.Take(1)
	fn.SeqTest(t, nums.Seq()).Is(6)

	var first fn.Opt[int]
	first, tail = tail.First()
	fn.OptTest(t, first).Is(8)

	nums, tail = tail.TakeWhile(func(i int) bool { return i < 20 })
	fn.SeqTest(t, nums.Seq()).Is(10, 12, 14, 16, 18)
}
