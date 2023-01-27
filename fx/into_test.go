package fx

import "testing"

func TestInto(t *testing.T) {
	vals := []int{1, 2, 3}
	sum := Into(0, func(res, e int) int {
		return res + e
	}, vals)

	if sum != 6 {
		t.Fatalf("bad result: %v", sum)
	}
}
