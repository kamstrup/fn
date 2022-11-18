package fn

import "testing"

func TestArraySeq(t *testing.T) {
	var arr Seq[int] = ArrayOfArgs(1, 2, 3)
	tarr := SeqTest(t, arr)
	tarr.LenIs(3)
	tarr.Is(1, 2, 3)

	arr = ArrayOf([]int{})
	tarr = SeqTest(t, arr)
	tarr.LenIs(0)
	tarr.IsEmpty()
}

func TestArraySeqTakeWhile(t *testing.T) {
	arr := ArrayOfArgs(1, 2, 3)
	head, tail := arr.TakeWhile(func(i int) bool { return i == 0 })
	SeqTest(t, head).IsEmpty()
	SeqTest(t, tail).Is(1, 2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 2 })
	SeqTest(t, head).Is(1)
	SeqTest(t, tail).Is(2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 3 })
	SeqTest(t, head).Is(1, 2)
	SeqTest(t, tail).Is(3)
	head, tail = arr.TakeWhile(func(_ int) bool { return true })
	SeqTest(t, head).Is(1, 2, 3)
	SeqTest(t, tail).IsEmpty()
}
