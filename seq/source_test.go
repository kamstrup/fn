package seq_test

import (
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestSource(t *testing.T) {
	nums, tail := seq.RangeFrom(0).
		Map(func(i int) int { return i * 2 }).
		Take(3)
	fntesting.TestOf(t, nums.Seq()).Is(0, 2, 4)

	nums, tail = tail.Take(1)
	fntesting.TestOf(t, nums.Seq()).Is(6)

	var first opt.Opt[int]
	first, tail = tail.First()
	fntesting.OptOf(t, first).Is(8)

	nums, tail = tail.TakeWhile(func(i int) bool { return i < 20 })
	fntesting.TestOf(t, nums.Seq()).Is(10, 12, 14, 16, 18)
}
