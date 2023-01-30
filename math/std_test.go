package fnmath

import (
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/opt"
)

func TestCollectSum(t *testing.T) {
	var arr fn.Seq[int] = fn.ArrayOfArgs(1, 2, 3)
	sum := fn.Into(0, Sum[int], arr)
	if sum.Must() != 6 {
		t.Errorf("expected sum 6: %d", sum)
	}

	sum = fn.Into(27, Sum[int], fn.SeqEmpty[int]())
	if !sum.Empty() || sum.Ok() {
		t.Errorf("expected empty sum: %v", sum)
	}
	if val, err := sum.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty sum: %v", sum)
	}
}

func TestCollectMinMax(t *testing.T) {
	arr := fn.ArrayOfArgs(1, 2, 3, 2, -1, 1)
	min := fn.Into(0, Min[int], arr)
	if min.Must() != -1 {
		t.Errorf("expected min -1: %d", min)
	}

	min = fn.Into(27, Min[int], fn.SeqEmpty[int]())
	if !min.Empty() || min.Ok() {
		t.Errorf("expected empty min: %v", min)
	}
	if val, err := min.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty min: %v", min)
	}

	max := fn.Into(0, Max[int], arr)
	if max.Must() != 3 {
		t.Errorf("expected max 3: %d", max)
	}

	max = fn.Into(27, Max[int], fn.SeqEmpty[int]())
	if !max.Empty() || max.Ok() {
		t.Errorf("expected empty max: %v", min)
	}
	if val, err := max.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty max: %v", max)
	}
}
