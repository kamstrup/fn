package fn

import "testing"

func TestSource(t *testing.T) {
	nums, tail := SourceOf(NumbersFrom(0)).Take(3)
	SeqTest(t, nums.Seq()).Is(0, 1, 2)

	nums, tail = tail.Take(1)
	SeqTest(t, nums.Seq()).Is(3)

	var first Opt[int]
	first, tail = tail.First()
	OptTest(t, first).Is(4)

	nums, tail = tail.TakeWhile(func(i int) bool { return i < 10 })
	SeqTest(t, nums.Seq()).Is(5, 6, 7, 8, 9)
}
