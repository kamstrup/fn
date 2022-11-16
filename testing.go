package fn

import (
	"testing"
)

func SeqTest[S comparable](t *testing.T, seq Seq[S]) TestSeq[S] {
	return TestSeq[S]{
		t:   t,
		seq: seq,
	}
}

type TestSeq[S comparable] struct {
	t   *testing.T
	seq Seq[S]
}

func (ts TestSeq[S]) IsExactly(ss ...S) {
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
