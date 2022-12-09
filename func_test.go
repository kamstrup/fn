package fn_test

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/kamstrup/fn"
)

func TestZeroes(t *testing.T) {
	fn.SeqTest(t, fn.MapOf(fn.ArrayOfArgs(-1, 0, 1, 10).Seq(), fn.IsZero[int])).Is(false, true, false, false)
	fn.SeqTest(t, fn.MapOf(fn.ArrayOfArgs(-1, 0, 1, 10).Seq(), fn.IsNonZero[int])).Is(true, false, true, true)
}

func TestCollectSum(t *testing.T) {
	var arr fn.Seq[int] = fn.ArrayOfArgs(1, 2, 3)
	sum := fn.Into(0, fn.Sum[int], arr)
	if sum != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = fn.Into(27, fn.Sum[int], fn.SeqEmpty[int]())
	if sum != 27 {
		t.Errorf("expected sum 6: %d", sum)
	}
}

func TestCollectMinMax(t *testing.T) {
	arr := fn.ArrayOfArgs(1, 2, 3, 2, -1, 1).Seq()
	min := fn.Into(0, fn.Min[int], arr)
	if min != -1 {
		t.Errorf("expected min -1: %d", min)
	}

	min = fn.Into(27, fn.Min[int], fn.SeqEmpty[int]())
	if min != 27 {
		t.Errorf("expected min 27: %d", min)
	}

	max := fn.Into(0, fn.Max[int], arr)
	if max != 3 {
		t.Errorf("expected max 3: %d", max)
	}

	max = fn.Into(27, fn.Max[int], fn.SeqEmpty[int]())
	if max != 27 {
		t.Errorf("expected max 27: %d", max)
	}
}

func TestCollectCount(t *testing.T) {
	arr := fn.ArrayOfArgs[int](1, 2, 3).Seq()

	count := fn.Into(0, fn.Count[int], arr)
	if count != 3 {
		t.Errorf("expected count 3: %d", count)
	}

	count = fn.Into(0, fn.Count[int], fn.SeqEmpty[int]())
	if count != 0 {
		t.Errorf("expected count 0: %d", count)
	}
}

func TestCollectAppend(t *testing.T) {
	arr := fn.ArrayOfArgs[int](1, 2, 3)
	cpy := fn.Into[int](nil, fn.Append[int], arr)
	exp := []int{1, 2, 3}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}

	cpy = fn.Into([]int{27}, fn.Append[int], fn.SeqEmpty[int]())
	exp = []int{27}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}
}

func TestCollectAssoc(t *testing.T) {
	oddNums := fn.ArrayOfArgs(1, 2, 3).
		Where(func(i int) bool { return i%2 == 1 })

	arr := fn.MapOf(oddNums, fn.TupleWithKey(func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	}))
	res := fn.Into(nil, fn.Assoc[string, int], arr)
	exp := map[string]int{
		"1": 1, "3": 3,
	}
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectSet(t *testing.T) {
	nums := fn.ArrayOfArgs(1, 2, 2, 3, 1).Seq()
	res := fn.Into(nil, fn.Set[int], nums)
	exp := map[int]struct{}{
		1: {}, 2: {}, 3: {},
	}
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectString(t *testing.T) {
	strs := fn.ArrayOfArgs("one", "two").Seq()
	res := fn.Into(nil, fn.StringBuilder, strs)
	exp := "onetwo"
	if exp != res.String() {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectGroupBy(t *testing.T) {
	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := fn.ZipOf[string, int](names, fn.NumbersFrom(0))
	res := fn.Into(nil, fn.GroupBy[string, int], tups)
	exp := map[string][]int{
		"bob":    {0, 2, 4},
		"alan":   {1, 5},
		"scotty": {3},
	}

	if !reflect.DeepEqual(exp, res) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectUpdateAssoc(t *testing.T) {
	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := fn.ZipOf[string, int](names, fn.Constant(1))
	res := fn.Into(nil, fn.UpdateAssoc[string, int](fn.Sum[int]), tups)
	exp := map[string]int{
		"bob":    3,
		"alan":   2,
		"scotty": 1,
	}

	if !reflect.DeepEqual(exp, res) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectUpdateArray(t *testing.T) {

	hellos := fn.ArrayOfArgs(
		fn.TupleOf(1, "hello"), fn.TupleOf(2, "hej"),
		fn.TupleOf(1, "world"), fn.TupleOf(2, "verden")).Seq()

	res := fn.Into(nil, fn.UpdateArray[int, string](func(old, new_ string) string {
		return strings.TrimSpace(old + " " + new_)
	}), hellos)

	exp := []string{
		"",
		"hello world",
		"hej verden",
	}

	if !reflect.DeepEqual(exp, res) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectErr(t *testing.T) {
	expErr := errors.New("the error")
	var nums fn.Seq[int]
	nums = fn.ArrayOfArgs(1, 2, 3)
	res, err := fn.IntoErr(0, func(into, n int) (int, error) {
		if into >= 2 {
			return 27, expErr
		}
		return into + 1, nil
	}, nums)

	if res != 27 {
		t.Errorf("expected 27, got %d", res)
	}

	if err != expErr {
		t.Errorf("did not get expected error: %v", expErr)
	}
}
