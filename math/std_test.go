package fnmath

import (
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

func TestCollectSum(t *testing.T) {
	var arr = seq.SliceOfArgs(1, 2, 3)
	sum := seq.Into(0, Sum[int], arr)
	if sum.Must() != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = seq.Into(27, Sum[int], seq.SeqEmpty[int]())
	if !sum.Empty() || sum.Ok() {
		t.Errorf("expected empty sum: %v", sum)
	}
	if val, err := sum.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty sum: %v", sum)
	}
}

func TestCollectMinMax(t *testing.T) {
	arr := seq.SliceOfArgs(1, 2, 3, 2, -1, 1)
	min := seq.Into(0, Min[int], arr)
	if min.Must() != -1 {
		t.Errorf("expected min -1: %d", min)
	}

	min = seq.Into(27, Min[int], seq.SeqEmpty[int]())
	if !min.Empty() || min.Ok() {
		t.Errorf("expected empty min: %v", min)
	}
	if val, err := min.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty min: %v", min)
	}

	max := seq.Into(0, Max[int], arr)
	if max.Must() != 3 {
		t.Errorf("expected max 3: %d", max)
	}

	max = seq.Into(27, Max[int], seq.SeqEmpty[int]())
	if !max.Empty() || max.Ok() {
		t.Errorf("expected empty max: %v", min)
	}
	if val, err := max.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty max: %v", max)
	}
}
