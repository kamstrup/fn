package slice

import (
	"reflect"
	"testing"
)

func TestMapping(t *testing.T) {
	even := Mapping([]int{1, 2, 3}, func(i int) int {
		return i * 2
	})
	if !reflect.DeepEqual(even, []int{2, 4, 6}) {
		t.Fatalf("bad results: %v", even)
	}
}

func TestMapIndex(t *testing.T) {
	results := MappingIndex([]int{1, 2, 3}, func(idx, n int) int {
		return n + idx
	})
	if !reflect.DeepEqual(results, []int{1, 3, 5}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestGen(t *testing.T) {
	results := Gen(3, func(idx int) int {
		return idx
	})
	if !reflect.DeepEqual(results, []int{0, 1, 2}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestCopy(t *testing.T) {
	orig := []int{1, 2, 3}
	results := Copy(orig)
	if !reflect.DeepEqual(results, []int{1, 2, 3}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestZero(t *testing.T) {
	orig := []int{1, 2, 3}
	results := Zero(orig)
	if !reflect.DeepEqual(results, []int{0, 0, 0}) || !reflect.DeepEqual(orig, []int{0, 0, 0}) {
		t.Fatalf("bad results: %v", results)
	}
}

func TestSortAsc(t *testing.T) {
	data := []int{1, 3, 2}
	SortAsc(data)
	if !reflect.DeepEqual(data, []int{1, 2, 3}) {
		t.Fatalf("bad results: %v", data)
	}
}

func TestSortDesc(t *testing.T) {
	data := []int{1, 3, 2}
	SortDesc(data)
	if !reflect.DeepEqual(data, []int{3, 2, 1}) {
		t.Fatalf("bad results: %v", data)
	}
}
