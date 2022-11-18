package fn

import (
	"reflect"
	"strconv"
	"testing"
)

func TestCollectSum(t *testing.T) {
	arr := ArrayOfArgs(1, 2, 3)

	sum := Collect[int](arr, Sum[int], 0)
	if sum != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = Collect(SeqEmpty[int](), Sum[int], 27)
	if sum != 27 {
		t.Errorf("expected sum 6: %d", sum)
	}
}

func TestCollectCount(t *testing.T) {
	arr := ArrayOfArgs[int](1, 2, 3)

	count := Collect[int](arr, Count[int], 0)
	if count != 3 {
		t.Errorf("expected count 3: %d", count)
	}

	count = Collect(SeqEmpty[int](), Count[int], 0)
	if count != 0 {
		t.Errorf("expected count 0: %d", count)
	}
}

func TestCollectAppend(t *testing.T) {
	arr := ArrayOfArgs[int](1, 2, 3)
	cpy := Collect[int](arr, Append[int], nil)
	exp := []int{1, 2, 3}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}

	cpy = Collect(SeqEmpty[int](), Append[int], []int{27})
	exp = []int{27}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}
}

func TestCollectAssoc(t *testing.T) {
	oddNums := ArrayOfArgs(1, 2, 3).
		Where(func(i int) bool { return i%2 == 1 })

	arr := SeqMap(oddNums, func(i int) Tuple[string, int] {
		return TupleOf(strconv.FormatInt(int64(i), 10), i)
	})
	res := Collect(arr, Assoc[string, int], map[string]int{})
	exp := map[string]int{
		"1": 1, "3": 3,
	}
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}
