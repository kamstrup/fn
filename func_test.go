package fn

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestZeroes(t *testing.T) {
	SeqTest(t, MapOf(ArrayOfArgs(-1, 0, 1, 10).Seq(), IsZero[int])).Is(false, true, false, false)
	SeqTest(t, MapOf(ArrayOfArgs(-1, 0, 1, 10).Seq(), IsNonZero[int])).Is(true, false, true, true)
}

func TestCollectSum(t *testing.T) {
	var arr Seq[int] = ArrayOfArgs(1, 2, 3)
	sum := Into(0, Sum[int], arr)
	if sum != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = Into(27, Sum[int], SeqEmpty[int]())
	if sum != 27 {
		t.Errorf("expected sum 6: %d", sum)
	}
}

func TestCollectCount(t *testing.T) {
	arr := ArrayOfArgs[int](1, 2, 3).Seq()

	count := Into(0, Count[int], arr)
	if count != 3 {
		t.Errorf("expected count 3: %d", count)
	}

	count = Into(0, Count[int], SeqEmpty[int]())
	if count != 0 {
		t.Errorf("expected count 0: %d", count)
	}
}

func TestCollectAppend(t *testing.T) {
	arr := ArrayOfArgs[int](1, 2, 3)
	cpy := Into[int](nil, Append[int], arr)
	exp := []int{1, 2, 3}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}

	cpy = Into([]int{27}, Append[int], SeqEmpty[int]())
	exp = []int{27}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}
}

func TestCollectAssoc(t *testing.T) {
	oddNums := ArrayOfArgs(1, 2, 3).
		Where(func(i int) bool { return i%2 == 1 })

	arr := MapOf(oddNums, TupleWithKey(func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	}))
	res := Into(nil, Assoc[string, int], arr)
	exp := map[string]int{
		"1": 1, "3": 3,
	}
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectSet(t *testing.T) {
	nums := ArrayOfArgs(1, 2, 2, 3, 1).Seq()
	res := Into(nil, Set[int], nums)
	exp := map[int]struct{}{
		1: {}, 2: {}, 3: {},
	}
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectString(t *testing.T) {
	strs := ArrayOfArgs("one", "two").Seq()
	res := Into(nil, StringBuilder, strs)
	exp := "onetwo"
	if exp != res.String() {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectGroupBy(t *testing.T) {
	names := ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := ZipOf[string, int](names, NumbersFrom(0))
	res := Into(nil, GroupBy[string, int], tups)
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
	names := ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := ZipOf[string, int](names, Constant(1))
	res := Into(nil, UpdateAssoc[string, int](Sum[int]), tups)
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

	hellos := ArrayOfArgs(
		TupleOf(1, "hello"), TupleOf(2, "hej"),
		TupleOf(1, "world"), TupleOf(2, "verden")).Seq()

	res := Into(nil, UpdateArray[int, string](func(old, new_ string) string {
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
	var nums Seq[int]
	nums = ArrayOfArgs(1, 2, 3)
	res, err := IntoErr(0, func(into, n int) (int, error) {
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
