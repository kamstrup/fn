package fn

import "testing"

func TestArraySeq(t *testing.T) {
	arr := ArrayOfArgs(1, 2, 3)
	tarr := SeqTest[int](t, arr)
	tarr.LenIs(3)
	tarr.Is(1, 2, 3)

	arr = ArrayOf([]int{})
	tarr = SeqTest[int](t, arr)
	tarr.LenIs(0)
	tarr.IsEmpty()
}

func TestArraySeqTakeWhile(t *testing.T) {
	arr := ArrayOfArgs(1, 2, 3)

	head, tail := arr.TakeWhile(func(i int) bool { return i == 0 })
	SeqTest[int](t, head).IsEmpty()
	SeqTest[int](t, tail).Is(1, 2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 2 })
	SeqTest[int](t, head).Is(1)
	SeqTest[int](t, tail).Is(2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 3 })
	SeqTest[int](t, head).Is(1, 2)
	SeqTest[int](t, tail).Is(3)
	head, tail = arr.TakeWhile(func(_ int) bool { return true })
	SeqTest[int](t, head).Is(1, 2, 3)
	SeqTest[int](t, tail).IsEmpty()
}
