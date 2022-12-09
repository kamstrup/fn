package fn

import "testing"

func TestSource(t *testing.T) {
	nums, tail := NumbersFrom(0).
		Map(func(i int) int { return i * 2 }).
		Take(3)
	SeqTest(t, nums.Seq()).Is(0, 2, 4)

	nums, tail = tail.Take(1)
	SeqTest(t, nums.Seq()).Is(6)

	var first Opt[int]
	first, tail = tail.First()
	OptTest(t, first).Is(8)

	nums, tail = tail.TakeWhile(func(i int) bool { return i < 20 })
	SeqTest(t, nums.Seq()).Is(10, 12, 14, 16, 18)
}
