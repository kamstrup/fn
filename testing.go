package fn

import (
	"fmt"
	"testing"
)

type TestSeq[S comparable] struct {
	t   *testing.T
	seq Seq[S]
}

type TestSeqSuite[S comparable] struct {
	t         *testing.T
	createSeq func() Seq[S]
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

func SeqTestSuite[S comparable](t *testing.T, createSeq func() Seq[S]) TestSeqSuite[S] {
	return TestSeqSuite[S]{
		t:         t,
		createSeq: createSeq,
	}
}

func (ts TestSeq[S]) LenIs(n int) {
	ts.t.Helper()

	if sz := ts.seq.Len(); sz != n {
		ts.t.Errorf("Seq len mismatch. Expected %d, found %d", n, sz)
	}
}

func (ts TestSeq[S]) Is(ss ...S) {
	ts.t.Helper()

	sz := ts.seq.Len()
	if sz != LenUnknown {
		if sz != len(ss) {
			ts.t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	count := 0
	ts.seq.ForEachIndex(func(i int, s S) {
		count++
		if i >= len(ss) {
			ts.t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, i, s)
		} else if ss[i] != s {
			ts.t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				i, ss[i], s)
		}
	})

	if sz != LenUnknown && sz != count {
		ts.t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, count)
	}
}

func (ts TestSeq[S]) IsEmpty() {
	ts.t.Helper()

	count := 0
	ts.seq.ForEach(func(s S) {
		count++
	})
	if count != 0 {
		ts.t.Errorf("Seq is not empty. Length %d", count)
	}
}

func (to TestOpt[S]) Is(s S) {
	to.t.Helper()

	val, err := to.opt.Return()
	if err != nil {
		to.t.Errorf("Option mismatch: %s", err)
	} else if val != s {
		to.t.Errorf("Option value mismatch. Expected %v, found %v", s, val)
	}
}

func (to TestOpt[S]) IsEmpty() {
	to.t.Helper()

	val, err := to.opt.Return()
	if err == nil {
		to.t.Errorf("option is not empty: %v", val)
	}
}

func (ts TestSeqSuite[S]) Is(ss ...S) {
	ts.t.Helper()

	ts.t.Run("ForEach", func(t *testing.T) {
		seqIsForEach(t, ts.createSeq(), ss)
	})

	ts.t.Run("ForEachIndex", func(t *testing.T) {
		seqIsForEachIndex(t, ts.createSeq(), ss)
	})

	ts.t.Run("Take", func(t *testing.T) {
		seqIsTake(t, ts.createSeq(), ss)
	})

	ts.t.Run("First", func(t *testing.T) {
		seqIsFirst(t, ts.createSeq(), ss)
	})
}

func (ts TestSeqSuite[S]) IsEmpty() {
	ts.t.Helper()

	count := 0
	ts.createSeq().ForEach(func(s S) {
		ts.t.Helper()
		count++
	})
	if count != 0 {
		ts.t.Errorf("Seq is not empty. Length %d", count)
	}
}

func seqIsForEach[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz := seq.Len()
	if sz != LenUnknown {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	i := 0
	seq.ForEach(func(s S) {
		t.Helper()

		if i >= len(ss) {
			t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, i, s)
		} else if ss[i] != s {
			t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				i, ss[i], s)
		}

		i++
	})

	if sz != LenUnknown && sz != i {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, i)
	}
}

func seqIsForEachIndex[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz := seq.Len()
	if sz != LenUnknown {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	count := 0
	seq.ForEachIndex(func(i int, s S) {
		t.Helper()

		count++
		if i >= len(ss) {
			t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, i, s)
		} else if ss[i] != s {
			t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				i, ss[i], s)
		}
	})

	if sz != LenUnknown && sz != count {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, count)
	}
}

func seqIsTake[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz := seq.Len()
	if sz != LenUnknown {
		t.Run("Len", func(t *testing.T) {
			if sz != len(ss) {
				t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
			}
		})
	}

	t.Run("0", func(t *testing.T) {
		head, tail := seq.Take(0)
		if head.Len() != 0 {
			t.Errorf("When calling Take(0) 'head' must be the empty array")
		}
		if sz != LenUnknown && sz != tail.Len() {
			t.Errorf("When calling Take(0) 'tail' must have the same length")
		}
	})

	// Ensure we can Take(n) for different n, and rebuild the exact ss
	for _, n := range []int{1, 2, 3, 100} {
		t.Run(fmt.Sprintf("%d", n), func(t *testing.T) {
			count := 0
			for head, tail := seq.Take(n); head.Len() != 0; head, tail = tail.Take(n) {
				for i := range head {
					if ss[count+i] != head[i] {
						t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
							count+i, ss[count], head[i])
					}
				}

				count += len(head)
				if count > len(ss) {
					t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
						len(ss)-1, count, head[0])
				}
			}

			if sz != LenUnknown && sz != count {
				t.Errorf("Number of elements in Take(%d) incorrect. Expected %d, found %d",
					n, sz, count)
			}
		})
	}
}

func seqIsFirst[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz := seq.Len()
	if sz != LenUnknown {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	var fst Opt[S]
	idx := 0
	for fst, seq = seq.First(); fst.Ok(); fst, seq = seq.First() {
		s := fst.val
		if idx >= len(ss) {
			t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, idx, s)
		} else if ss[idx] != s {
			t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				idx, ss[idx], s)
		}
		idx++
	}

	if sz != LenUnknown && sz != idx {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, idx)
	}
}
