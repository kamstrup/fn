package slice

import "testing"

func TestInto(t *testing.T) {
	vals := []int{1, 2, 3}
	sum := Reduce(func(res, e int) int {
		return res + e
	}, 0, vals)

	if sum != 6 {
		t.Fatalf("bad result: %v", sum)
	}
}
