package fx

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMapKeys(t *testing.T) {
	m := map[string]int{
		"one": 1, "two": 2,
	}
	keys := SortSliceAsc(Keys(m))
	if !reflect.DeepEqual([]string{"one", "two"}, keys) {
		t.Fatalf("bad results: %v", keys)
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{
		"one": 1, "two": 2,
	}
	vals := SortSliceAsc(Values(m))
	if !reflect.DeepEqual([]int{1, 2}, vals) {
		t.Fatalf("bad results: %v", vals)
	}
}

func TestMapAssoc(t *testing.T) {
	m := map[string]int{
		"one": 1, "two": 2,
	}

	m2 := MapAssoc(m, func(k string, v int) (string, string) {
		return k, strconv.FormatInt(int64(v), 10)
	})
	if !reflect.DeepEqual(m2, map[string]string{
		"one": "1", "two": "2",
	}) {
		t.Fatalf("bad results: %v", m2)
	}

	// with key collisions we should end up w a smaller map
	m3 := MapAssoc(m, func(k string, _ int) (int, string) {
		return len(k), k
	})
	res1 := map[int]string{3: "one"} // since key order is random we can get res1 or res2
	res2 := map[int]string{3: "two"}
	if !(reflect.DeepEqual(m3, res1) || reflect.DeepEqual(m3, res2)) {
		t.Fatalf("bad results: %v", m3)
	}
}
