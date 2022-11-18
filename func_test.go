package fn

import (
	"reflect"
	"testing"
)

func TestCollectSum(t *testing.T) {
	arr := ArrayOfArgs[int](1, 2, 3)

	sum := Collect[int](arr, Sum[int], 0)
	if sum != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = Collect[int](SeqEmpty[int](), Sum[int], 27)
	if sum != 27 {
		t.Errorf("expected sum 6: %d", sum)
	}
}

func TestCollectAppend(t *testing.T) {
	arr := ArrayOfArgs[int](1, 2, 3)
	cpy := Collect[int](arr, Append[int], nil)
	exp := []int{1, 2, 3}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}

	cpy = Collect[int](SeqEmpty[int](), Append[int], []int{27})
	exp = []int{27}
	if !reflect.DeepEqual(cpy, exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}
}
