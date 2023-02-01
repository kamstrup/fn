package slice

import (
	"reflect"
	"testing"
)

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
