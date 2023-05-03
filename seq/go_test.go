package seq_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestGo(t *testing.T) {
	numStrs := seq.Go(seq.SliceOfArgs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9), 5, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})
	numStrSorted := numStrs.ToSlice().Sort(seq.OrderAsc[string])
	fntesting.TestOf(t, numStrSorted.Seq()).Is("0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
}

func TestGoErrorTail(t *testing.T) {
	// The intent of this test is to assert that on errors the seq returned from ForEach is also an error seq.
	theError := errors.New("the error")
	s := seq.ConcatOf(seq.SliceOfArgs(1, 2, 3), seq.ErrorOf[int](theError))
	numStrs := seq.Go(s, 3, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})

	arr := make([]string, 3)
	res := numStrs.ForEachIndex(func(i int, s string) {
		arr[i] = s
	})

	err := res.Error()
	if err != theError {
		t.Fatalf("expected %q, got: %q", theError, err)
	}

	// sort so we can presume an order
	arr = seq.SliceAs(arr).Sort(seq.OrderAsc[string])
	if !reflect.DeepEqual(arr, []string{"1", "2", "3"}) {
		t.Fatalf("bad result: %v", arr)
	}
}

func TestGoErrorOpt(t *testing.T) {
	// The intent of this test is to assert that if the input ends with an error,
	// and the output type is Opt, then the last element in the output is an error opt.
	theError := errors.New("the error")
	ints := seq.ConcatOf(seq.SliceOfArgs(1, 2, 3), seq.ErrorOf[int](theError))
	optStrs := seq.Go(ints, 2, func(i int) opt.Opt[string] {
		return opt.Of(strconv.FormatInt(int64(i), 10))
	})

	arr := make([]string, 3)
	goTail := optStrs.ForEachIndex(func(i int, so opt.Opt[string]) {
		s, err := so.Return()
		if err != nil {
			if i != 3 {
				t.Fatalf("error must be last element: %d", i)
			}
			if err != theError {
				t.Fatalf("expected %q, got: %q", theError, err)
			}
		} else {
			arr[i] = s
		}
	})

	// sort so we can presume an order
	arr = seq.SliceAs(arr).Sort(seq.OrderAsc[string])
	if !reflect.DeepEqual(arr, []string{"1", "2", "3"}) {
		t.Fatalf("bad result: %v", arr)
	}

	if goTail.Error() != theError {
		t.Fatalf("tail from Go.ForEach does not have correct error: %v", goTail)
	}
}
