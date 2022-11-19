package fn

import "testing"

func TestWhile(t *testing.T) {
	nums := SourceOf(NumbersFrom(0)).
		While(func(i int) bool { return i < 10 }).
		Array().Seq()
	SeqTest(t, nums).Is(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
}
