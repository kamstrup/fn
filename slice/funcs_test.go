package slice

import (
	"reflect"
	"testing"
)

func TestUniq(t *testing.T) {
	nums := Uniq([]string{"one", "two", "one", "three"})
	if !reflect.DeepEqual(nums, []string{"one", "two", "three"}) {
		t.Fatalf("bad result: %v", nums)
	}
}

func TestGroupBy(t *testing.T) {
	words := []string{"one", "two", "three", "four"}
	byFirstByte := GroupBy(words, func(word string) uint8 {
		return word[0]
	})
	if !reflect.DeepEqual(byFirstByte, map[uint8][]string{
		'o': []string{"one"},
		't': []string{"two", "three"},
		'f': []string{"four"},
	}) {
		t.Fatalf("bad result: %v", byFirstByte)
	}
}

func TestReduce(t *testing.T) {
	vals := []int{1, 2, 3}
	sum := Reduce(func(res, e int) int {
		return res + e
	}, 0, vals)

	if sum != 6 {
		t.Fatalf("bad result: %v", sum)
	}
}
