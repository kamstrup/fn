package slice

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMapKeys(t *testing.T) {
	m := map[string]int{
		"one": 1, "two": 2,
	}
	keys := SortAsc(Keys(m))
	if !reflect.DeepEqual([]string{"one", "two"}, keys) {
		t.Fatalf("bad results: %v", keys)
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{
		"one": 1, "two": 2,
	}
	vals := SortAsc(Values(m))
	if !reflect.DeepEqual([]int{1, 2}, vals) {
		t.Fatalf("bad results: %v", vals)
	}
}

func TestToMap(t *testing.T) {
	s := []int{1, 2, 3}
	m := ToMap(s, func(i int) (string, int8) {
		return strconv.FormatInt(int64(i), 10), int8(i)
	})
	if !reflect.DeepEqual(m, map[string]int8{
		"1": 1, "2": 2, "3": 3,
	}) {
		t.Fatalf("bad result: %v", m)
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]int8{
		"1": 1, "2": 2, "3": 3,
	}
	s := FromMap(m, func(k string, i int8) int {
		return int(i)
	})
	SortAsc(s)
	if !reflect.DeepEqual([]int{1, 2, 3}, s) {
		t.Fatalf("bad result: %v", s)
	}
}
