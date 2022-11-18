package fn

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestCollectSum(t *testing.T) {
	var arr Seq[int] = ArrayOfArgs(1, 2, 3)
	sum := Collect(arr, Sum[int], 0)
	if sum != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = Collect(SeqEmpty[int](), Sum[int], 27)
	if sum != 27 {
		t.Errorf("expected sum 6: %d", sum)
	}
}

func TestCollectCount(t *testing.T) {
	var arr Seq[int]
	arr = ArrayOfArgs[int](1, 2, 3)

	count := Collect(arr, Count[int], 0)
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

func TestCollectErr(t *testing.T) {
	expErr := errors.New("the error")
	var nums Seq[int]
	nums = ArrayOfArgs(1, 2, 3)
	res, err := CollectErr(nums,
		func(into, n int) (int, error) {
			if into >= 2 {
				return 27, expErr
			}
			return into + 1, nil
		}, 0)

	if res != 27 {
		t.Errorf("expected 27, got %d", res)
	}

	if err != expErr {
		t.Errorf("did not get expected error: %v", expErr)
	}
}
