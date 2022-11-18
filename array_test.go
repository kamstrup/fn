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
	var (
		arr, head Array[int]
		tail      Seq[int]
	)
	arr = ArrayOfArgs(1, 2, 3)
	head, tail = arr.TakeWhile(func(i int) bool { return i == 0 })
	SeqTest(t, head.Seq()).IsEmpty()
	SeqTest(t, tail).Is(1, 2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 2 })
	SeqTest(t, head.Seq()).Is(1)
	SeqTest(t, tail).Is(2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 3 })
	SeqTest(t, head.Seq()).Is(1, 2)
	SeqTest(t, tail).Is(3)
	head, tail = arr.TakeWhile(func(_ int) bool { return true })
	SeqTest(t, head.Seq()).Is(1, 2, 3)
	SeqTest(t, tail).IsEmpty()
}

func TestArraySort(t *testing.T) {
	arr := ArrayOfArgs(1, 2, 3, 4).Sort(OrderDesc[int]).Seq()
	SeqTest(t, arr).Is(4, 3, 2, 1)

	arr = ArrayOfArgs(2, 1, 3).Sort(OrderAsc[int]).Seq()
	SeqTest(t, arr).Is(1, 2, 3)

	arrTup := ArrayOfArgs(TupleOf(1, 1), TupleOf(2, 2)).Sort(OrderTupleDesc[int, int]).Seq()
	SeqTest(t, arrTup).Is(TupleOf(2, 2), TupleOf(1, 1))

	arrTup = ArrayOfArgs(TupleOf(2, 2), TupleOf(1, 1)).Sort(OrderTupleAsc[int, int]).Seq()
	SeqTest(t, arrTup).Is(TupleOf(1, 1), TupleOf(2, 2))
}

func TestArrayReverse(t *testing.T) {
	arr := ArrayOfArgs(1, 2, 3, 4).Reverse()
	SeqTest(t, arr).Is(4, 3, 2, 1)

	arr = ArrayOfArgs(1, 2, 3).Reverse()
	SeqTest(t, arr).Is(3, 2, 1)

	arr = ArrayOfArgs(1).Reverse()
	SeqTest(t, arr).Is(1)

	arr = ArrayOfArgs[int]().Reverse()
	SeqTest(t, arr).IsEmpty()
}
