package fn

import (
	"testing"
)

type TestSeq[S comparable] struct {
	t   *testing.T
	seq Seq[S]
}

type TestOpt[S comparable] struct {
	t   *testing.T
	opt Opt[S]
}

func SeqTest[S comparable](t *testing.T, seq Seq[S]) TestSeq[S] {
	return TestSeq[S]{
		t:   t,
		seq: seq,
	}
}

func OptTest[S comparable](t *testing.T, opt Opt[S]) TestOpt[S] {
	return TestOpt[S]{
		t:   t,
		opt: opt,
	}
}

func (ts TestSeq[S]) Is(ss ...S) {
	if sz := ts.seq.Len(); sz != LenUnknown {
		if sz != len(ss) {
			ts.t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	ts.seq.ForEachIndex(func(i int, s S) {
		if ss[i] != s {
			ts.t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				i, ss[i], s)
		}
	})
}

func (ts TestSeq[S]) IsEmpty() {
	count := 0
	ts.seq.ForEach(func(s S) {
		count++
	})
	if count != 0 {
		ts.t.Errorf("Seq is not empty. Length %d", count)
	}
}

func (to TestOpt[S]) Is(s S) {
	val, err := to.opt.Return()
	if err != nil {
		to.t.Errorf("Option mismatch: %s", err)
	} else if val != s {
		to.t.Errorf("Option value mismatch. Expected %v, found %v", s, val)
	}
}

func (to TestOpt[S]) IsEmpty() {
	val, err := to.opt.Return()
	if err == nil {
		to.t.Errorf("option is not empty: %v", val)
	}
}
