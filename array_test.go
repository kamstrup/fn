package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
)

func TestArraySeq(t *testing.T) {
	var arr fn.Seq[int] = fn.ArrayOfArgs(1, 2, 3)
	tarr := fn.SeqTest(t, arr)
	tarr.LenIs(3)
	tarr.Is(1, 2, 3)

	arr = fn.ArrayOf([]int{})
	tarr = fn.SeqTest(t, arr)
	tarr.LenIs(0)
	tarr.IsEmpty()
}

func TestArraySeqTakeWhile(t *testing.T) {
	var (
		arr, head fn.Array[int]
		tail      fn.Seq[int]
	)
	arr = fn.ArrayOfArgs(1, 2, 3)
	head, tail = arr.TakeWhile(func(i int) bool { return i == 0 })
	fn.SeqTest(t, head.Seq()).IsEmpty()
	fn.SeqTest(t, tail).Is(1, 2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 2 })
	fn.SeqTest(t, head.Seq()).Is(1)
	fn.SeqTest(t, tail).Is(2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 3 })
	fn.SeqTest(t, head.Seq()).Is(1, 2)
	fn.SeqTest(t, tail).Is(3)
	head, tail = arr.TakeWhile(func(_ int) bool { return true })
	fn.SeqTest(t, head.Seq()).Is(1, 2, 3)
	fn.SeqTest(t, tail).IsEmpty()
}

func TestArraySort(t *testing.T) {
	arr := fn.ArrayOfArgs(1, 2, 3, 4).Sort(fn.OrderDesc[int]).Seq()
	fn.SeqTest(t, arr).Is(4, 3, 2, 1)

	arr = fn.ArrayOfArgs(2, 1, 3).Sort(fn.OrderAsc[int]).Seq()
	fn.SeqTest(t, arr).Is(1, 2, 3)

	arrTup := fn.ArrayOfArgs(fn.TupleOf(1, 1), fn.TupleOf(2, 2)).Sort(fn.OrderTupleDesc[int, int]).Seq()
	fn.SeqTest(t, arrTup).Is(fn.TupleOf(2, 2), fn.TupleOf(1, 1))

	arrTup = fn.ArrayOfArgs(fn.TupleOf(2, 2), fn.TupleOf(1, 1)).Sort(fn.OrderTupleAsc[int, int]).Seq()
	fn.SeqTest(t, arrTup).Is(fn.TupleOf(1, 1), fn.TupleOf(2, 2))
}

func TestArrayReverse(t *testing.T) {
	arr := fn.ArrayOfArgs(1, 2, 3, 4).Reverse()
	fn.SeqTest(t, arr).Is(4, 3, 2, 1)

	arr = fn.ArrayOfArgs(1, 2, 3).Reverse()
	fn.SeqTest(t, arr).Is(3, 2, 1)

	arr = fn.ArrayOfArgs(1).Reverse()
	fn.SeqTest(t, arr).Is(1)

	arr = fn.ArrayOfArgs[int]().Reverse()
	fn.SeqTest(t, arr).IsEmpty()
}

func TestArraySuite(t *testing.T) {
	createArr := func() fn.Seq[int] {
		return fn.ArrayOfArgs(1, 2, 3, 4)
	}
	fn.SeqTestSuite(t, createArr).Is(1, 2, 3, 4)
}
