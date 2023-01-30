package seq_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestArraySeq(t *testing.T) {
	// var arr seq.Seq[int] = seq.SliceOfArgs(1, 2, 3)
	arr := seq.SliceOfArgs(1, 2, 3)
	tarr := fntesting.TestOf(t, arr)
	tarr.LenIs(3)
	tarr.Is(1, 2, 3)

	arr = seq.SliceOf([]int{})
	tarr = fntesting.TestOf(t, arr)
	tarr.LenIs(0)
	tarr.IsEmpty()
}

func TestArraySeqTakeWhile(t *testing.T) {
	var (
		arr, head seq.Slice[int]
		tail      seq.Seq[int]
	)
	arr = seq.SliceAsArgs(1, 2, 3)
	head, tail = arr.TakeWhile(func(i int) bool { return i == 0 })
	fntesting.TestOf(t, head.Seq()).IsEmpty()
	fntesting.TestOf(t, tail).Is(1, 2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 2 })
	fntesting.TestOf(t, head.Seq()).Is(1)
	fntesting.TestOf(t, tail).Is(2, 3)

	head, tail = arr.TakeWhile(func(i int) bool { return i < 3 })
	fntesting.TestOf(t, head.Seq()).Is(1, 2)
	fntesting.TestOf(t, tail).Is(3)
	head, tail = arr.TakeWhile(func(_ int) bool { return true })
	fntesting.TestOf(t, head.Seq()).Is(1, 2, 3)
	fntesting.TestOf(t, tail).IsEmpty()
}

func TestArraySort(t *testing.T) {
	arr := seq.SliceAsArgs(1, 2, 3, 4).Sort(seq.OrderDesc[int]).Seq()
	fntesting.TestOf(t, arr).Is(4, 3, 2, 1)

	arr = seq.SliceAsArgs(2, 1, 3).Sort(seq.OrderAsc[int]).Seq()
	fntesting.TestOf(t, arr).Is(1, 2, 3)

	arrTup := seq.SliceAsArgs(seq.TupleOf(1, 1), seq.TupleOf(2, 2)).Sort(seq.OrderTupleDesc[int, int]).Seq()
	fntesting.TestOf(t, arrTup).Is(seq.TupleOf(2, 2), seq.TupleOf(1, 1))

	arrTup = seq.SliceAsArgs(seq.TupleOf(2, 2), seq.TupleOf(1, 1)).Sort(seq.OrderTupleAsc[int, int]).Seq()
	fntesting.TestOf(t, arrTup).Is(seq.TupleOf(1, 1), seq.TupleOf(2, 2))
}

func TestArrayReverse(t *testing.T) {
	arr := seq.SliceAsArgs(1, 2, 3, 4).Reverse()
	fntesting.TestOf(t, arr).Is(4, 3, 2, 1)

	arr = seq.SliceAsArgs(1, 2, 3).Reverse()
	fntesting.TestOf(t, arr).Is(3, 2, 1)

	arr = seq.SliceAsArgs(1).Reverse()
	fntesting.TestOf(t, arr).Is(1)

	arr = seq.SliceAsArgs[int]().Reverse()
	fntesting.TestOf(t, arr).IsEmpty()
}

func TestArraySuite(t *testing.T) {
	createArr := func() seq.Seq[int] {
		return seq.SliceOfArgs(1, 2, 3, 4)
	}
	fntesting.SuiteOf(t, createArr).Is(1, 2, 3, 4)
}

func TestArrayError(t *testing.T) {
	theError := errors.New("the error")
	arr := seq.SliceOfArgs(seq.ErrorOf[int](theError))

	if err := seq.Error(arr); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
