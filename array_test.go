package fn_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestArraySeq(t *testing.T) {
	// var arr fn.Seq[int] = fn.SliceOfArgs(1, 2, 3)
	arr := fn.SliceOfArgs(1, 2, 3)
	tarr := fntesting.TestOf(t, arr)
	tarr.LenIs(3)
	tarr.Is(1, 2, 3)

	arr = fn.SliceOf([]int{})
	tarr = fntesting.TestOf(t, arr)
	tarr.LenIs(0)
	tarr.IsEmpty()
}

func TestArraySeqTakeWhile(t *testing.T) {
	var (
		arr, head fn.Slice[int]
		tail      fn.Seq[int]
	)
	arr = fn.SliceAsArgs(1, 2, 3)
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
	arr := fn.SliceAsArgs(1, 2, 3, 4).Sort(fn.OrderDesc[int]).Seq()
	fntesting.TestOf(t, arr).Is(4, 3, 2, 1)

	arr = fn.SliceAsArgs(2, 1, 3).Sort(fn.OrderAsc[int]).Seq()
	fntesting.TestOf(t, arr).Is(1, 2, 3)

	arrTup := fn.SliceAsArgs(fn.TupleOf(1, 1), fn.TupleOf(2, 2)).Sort(fn.OrderTupleDesc[int, int]).Seq()
	fntesting.TestOf(t, arrTup).Is(fn.TupleOf(2, 2), fn.TupleOf(1, 1))

	arrTup = fn.SliceAsArgs(fn.TupleOf(2, 2), fn.TupleOf(1, 1)).Sort(fn.OrderTupleAsc[int, int]).Seq()
	fntesting.TestOf(t, arrTup).Is(fn.TupleOf(1, 1), fn.TupleOf(2, 2))
}

func TestArrayReverse(t *testing.T) {
	arr := fn.SliceAsArgs(1, 2, 3, 4).Reverse()
	fntesting.TestOf(t, arr).Is(4, 3, 2, 1)

	arr = fn.SliceAsArgs(1, 2, 3).Reverse()
	fntesting.TestOf(t, arr).Is(3, 2, 1)

	arr = fn.SliceAsArgs(1).Reverse()
	fntesting.TestOf(t, arr).Is(1)

	arr = fn.SliceAsArgs[int]().Reverse()
	fntesting.TestOf(t, arr).IsEmpty()
}

func TestArraySuite(t *testing.T) {
	createArr := func() fn.Seq[int] {
		return fn.SliceOfArgs(1, 2, 3, 4)
	}
	fntesting.SuiteOf(t, createArr).Is(1, 2, 3, 4)
}

func TestArrayError(t *testing.T) {
	theError := errors.New("the error")
	arr := fn.SliceOfArgs(fn.ErrorOf[int](theError))

	if err := fn.Error(arr); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
