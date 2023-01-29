package fn_test

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/kamstrup/fn"
	fnmath "github.com/kamstrup/fn/math"
	"github.com/kamstrup/fn/testing"
)

func TestZeroes(t *testing.T) {
	fntesting.TestOf(t, fn.MapOf(fn.ArrayOfArgs(-1, 0, 1, 10), fn.IsZero[int])).Is(false, true, false, false)
	fntesting.TestOf(t, fn.MapOf(fn.ArrayOfArgs(-1, 0, 1, 10), fn.IsNonZero[int])).Is(true, false, true, true)
}

func TestCollectCount(t *testing.T) {
	arr := fn.ArrayOfArgs[int](1, 2, 3)

	count := fn.Into(0, fn.Count[int], arr)
	if count.Must() != 3 {
		t.Errorf("expected count 3: %d", count)
	}

	count = fn.Into(0, fn.Count[int], fn.SeqEmpty[int]())
	if !count.Empty() || count.Ok() {
		t.Errorf("expected empty count: %v", count)
	}
	if val, err := count.Return(); val == 27 || err != fn.ErrEmpty {
		t.Errorf("expected empty count: %v", count)
	}
}

func TestCollectAppend(t *testing.T) {
	arr := fn.ArrayOfArgs[int](1, 2, 3)
	cpy := fn.Into(nil, fn.Append[int], arr)
	exp := []int{1, 2, 3}
	if !reflect.DeepEqual(cpy.Must(), exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}

	ints := fn.Into([]int{27}, fn.Append[int], fn.SeqEmpty[int]())
	if !ints.Empty() || ints.Ok() {
		t.Errorf("expected empty min: %v", ints)
	}
	if val, err := ints.Return(); len(val) != 0 || err != fn.ErrEmpty {
		t.Errorf("expected empty ints: %v", ints)
	}
}

func TestCollectAssoc(t *testing.T) {
	oddNums := fn.ArrayOfArgs(1, 2, 3).
		Where(func(i int) bool { return i%2 == 1 })

	arr := fn.MapOf(oddNums, fn.TupleWithKey(func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	}))
	res := fn.Into(nil, fn.MakeAssoc[string, int], arr)
	exp := map[string]int{
		"1": 1, "3": 3,
	}
	if !reflect.DeepEqual(res.Must(), exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectSet(t *testing.T) {
	nums := fn.ArrayOfArgs(1, 2, 2, 3, 1)
	res := fn.Into(nil, fn.MakeSet[int], nums)
	exp := map[int]struct{}{
		1: {}, 2: {}, 3: {},
	}
	if !reflect.DeepEqual(res.Must(), exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectString(t *testing.T) {
	strs := fn.ArrayOfArgs("one", "two")
	res := fn.Into(nil, fn.MakeString, strs)
	exp := "onetwo"
	if exp != res.Must().String() {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectGroupBy(t *testing.T) {
	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := fn.ZipOf[string, int](names, fn.RangeFrom(0))
	res := fn.Into(nil, fn.GroupBy[string, int], tups)
	exp := map[string][]int{
		"bob":    {0, 2, 4},
		"alan":   {1, 5},
		"scotty": {3},
	}

	if !reflect.DeepEqual(exp, res.Must()) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectUpdateAssoc(t *testing.T) {
	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := fn.ZipOf[string, int](names, fn.Constant(1))
	res := fn.Into(nil, fn.UpdateAssoc[string, int](fnmath.Sum[int]), tups)
	exp := map[string]int{
		"bob":    3,
		"alan":   2,
		"scotty": 1,
	}

	if !reflect.DeepEqual(exp, res.Must()) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectUpdateArray(t *testing.T) {

	hellos := fn.ArrayOfArgs(
		fn.TupleOf(1, "hello"), fn.TupleOf(2, "hej"),
		fn.TupleOf(1, "world"), fn.TupleOf(2, "verden"))

	res := fn.Into(nil, fn.UpdateArray[int, string](func(old, new_ string) string {
		return strings.TrimSpace(old + " " + new_)
	}), hellos)

	exp := []string{
		"",
		"hello world",
		"hej verden",
	}

	if !reflect.DeepEqual(exp, res.Must()) {
		t.Errorf("expected %v, got %v", exp, res.Must())
	}
}

func TestCollectError(t *testing.T) {
	theError := errors.New("the error")
	errSeq := fn.ErrorOf[int](theError)
	res := fn.Into(0, fnmath.Sum[int], errSeq)

	if res.Error() != theError || res.Ok() {
		t.Errorf("expected 'the error': %v", res)
	}
}

func TestPredicates(t *testing.T) {
	if fn.LessThanZero(1) {
		t.Errorf("1 is not < 0")
	}
	if !fn.LessThanZero(-1) {
		t.Errorf("-1 is < 0")
	}

	if !fn.GreaterThanZero(1) {
		t.Errorf("1 is > 0")
	}
	if fn.GreaterThanZero(-1) {
		t.Errorf("-1 is not > 0")
	}

	if fn.Is("hello")("hej") {
		t.Errorf("hello != hej")
	}
	if !fn.IsNot("hello")("hej") {
		t.Errorf("hello != hej")
	}

	if fn.IsZero("hello") {
		t.Errorf("hello is non-zero")
	}
	if !fn.IsZero("") {
		t.Errorf("\"\" should be zero")
	}
	if !fn.IsNonZero("hello") {
		t.Errorf("hello is non-zero")
	}
	if fn.IsNonZero("") {
		t.Errorf("\"\" should be zero")
	}
}

func TestAny(t *testing.T) {
	if fn.Any(fn.ArrayOfArgs(1, 2, 3), fn.IsZero[int]) {
		t.Fatal("should not find zero")
	}

	if !fn.Any(fn.ArrayOfArgs(0, 1), fn.IsNonZero[int]) {
		t.Fatal("should find non-zero")
	}
}

func TestAll(t *testing.T) {
	if fn.All(fn.ArrayOfArgs(0, 0, 1), fn.IsZero[int]) {
		t.Fatal("should find no.zero")
	}

	if !fn.All(fn.ArrayOfArgs(0, 0), fn.IsZero[int]) {
		t.Fatal("should be all zeroes")
	}
}

func TestLast(t *testing.T) {
	o := fn.Last(fn.ArrayOfArgs(0, 0, 1))
	if o.Error() != nil {
		t.Fatal("should not error", o.Error())
	}
	if o.Must() != 1 {
		t.Fatal("should be 1", o.Must())
	}

	o = fn.Last(fn.ArrayOfArgs(0))
	if o.Error() != nil {
		t.Fatal("should not error", o.Error())
	}
	if o.Must() != 0 {
		t.Fatal("should be 0", o.Must())
	}

	o = fn.Last(fn.SeqEmpty[int]())
	if !o.Empty() {
		t.Fatal("should be empty", o.Must())
	}

	theError := errors.New("the error")
	o = fn.Last(fn.ErrorOf[int](theError))
	if !o.Empty() || o.Ok() {
		t.Fatal("should be empty", o.Must())
	}
	if o.Error() != theError {
		t.Fatal("should be 'the error'", o.Error())
	}
}
