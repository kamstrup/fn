package fnmath

import (
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

func TestCollectSum(t *testing.T) {
	var arr = seq.SliceOfArgs(1, 2, 3)
	sum := seq.Reduce(Sum[int], 0, arr)
	if sum.Must() != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = seq.Reduce(Sum[int], 27, seq.Empty[int]())
	if !sum.Empty() || sum.Ok() {
		t.Errorf("expected empty sum: %v", sum)
	}
	if val, err := sum.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty sum: %v", sum)
	}
}

func TestCollectMinMax(t *testing.T) {
	arr := seq.SliceOfArgs(1, 2, 3, 2, -1, 1)
	min := seq.Reduce(Min[int], 0, arr)
	if min.Must() != -1 {
		t.Errorf("expected min -1: %d", min)
	}

	min = seq.Reduce(Min[int], 27, seq.Empty[int]())
	if !min.Empty() || min.Ok() {
		t.Errorf("expected empty min: %v", min)
	}
	if val, err := min.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty min: %v", min)
	}

	max := seq.Reduce(Max[int], 0, arr)
	if max.Must() != 3 {
		t.Errorf("expected max 3: %d", max)
	}

	max = seq.Reduce(Max[int], 27, seq.Empty[int]())
	if !max.Empty() || max.Ok() {
		t.Errorf("expected empty max: %v", min)
	}
	if val, err := max.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty max: %v", max)
	}
}
