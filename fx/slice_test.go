package fx

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSliceMap(t *testing.T) {
	even := MapSlice([]int{1, 2, 3}, func(i int) int {
		return i * 2
	})
	if !reflect.DeepEqual(even, []int{2, 4, 6}) {
		t.Fatalf("bad results: %v", even)
	}
}

func TestSliceMapIndex(t *testing.T) {
	results := MapSliceIndex([]int{1, 2, 3}, func(idx, n int) int {
		return n + idx
	})
	if !reflect.DeepEqual(results, []int{1, 3, 5}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestSliceGen(t *testing.T) {
	results := GenSlice(3, func(idx int) int {
		return idx
	})
	if !reflect.DeepEqual(results, []int{0, 1, 2}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestSliceCopy(t *testing.T) {
	orig := []int{1, 2, 3}
	results := CopySlice(orig)
	if !reflect.DeepEqual(results, []int{1, 2, 3}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestSliceZero(t *testing.T) {
	orig := []int{1, 2, 3}
	results := ZeroSlice(orig)
	if !reflect.DeepEqual(results, []int{0, 0, 0}) || !reflect.DeepEqual(orig, []int{0, 0, 0}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestSliceSortAsc(t *testing.T) {
	data := []int{1, 3, 2}
	SortSliceAsc(data)
	if !reflect.DeepEqual(data, []int{1, 2, 3}) {
		t.Fatalf("bad results: %v", data)
	}
}

func TestSliceSortDesc(t *testing.T) {
	data := []int{1, 3, 2}
	SortSliceDesc(data)
	if !reflect.DeepEqual(data, []int{3, 2, 1}) {
		t.Fatalf("bad results: %v", data)
	}
}

func TestSliceAssoc(t *testing.T) {
	s := []int{1, 2, 3}
	m := AssocSlice(s, func(i int) (string, int8) {
		return strconv.FormatInt(int64(i), 10), int8(i)
	})
	if !reflect.DeepEqual(m, map[string]int8{
		"1": 1, "2": 2, "3": 3,
	}) {
		t.Fatalf("bad result: %v", m)
	}
}
