package fn

import (
	"fmt"
	"reflect"
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

	if sz, _ := ts.seq.Len(); sz != n {
		ts.t.Errorf("Seq len mismatch. Expected %d, found %d", n, sz)
	}
}

func (ts TestSeq[S]) Is(ss ...S) {
	ts.t.Helper()

	sz, lenOk := ts.seq.Len()
	if lenOk {
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

	if lenOk && sz != count {
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
		seqIsTake(t, ts.createSeq, ss)
	})

	ts.t.Run("TakeWhile", func(t *testing.T) {
		seqIsTakeWhile(t, ts.createSeq, ss)
	})

	ts.t.Run("Skip", func(t *testing.T) {
		seqSkip(t, ts.createSeq, ss)
	})

	ts.t.Run("Where", func(t *testing.T) {
		seqIsWhere(t, ts.createSeq, ss)
	})
	// TODO: test While(pred)

	ts.t.Run("First", func(t *testing.T) {
		seqIsFirst(t, ts.createSeq(), ss)
	})

	ts.t.Run("All", func(t *testing.T) {
		seqIsAll(t, ts.createSeq, ss)
	})

	ts.t.Run("Any", func(t *testing.T) {
		seqIsAny(t, ts.createSeq, ss)
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

	if arr := ts.createSeq().Array(); len(arr) != 0 {
		ts.t.Errorf("Seq.Array is not empty. Length %d", len(arr))
	}

	// TODO: more checks for emptiness!
}

func seqIsForEach[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz, lenOk := seq.Len()
	if lenOk {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	i := 0
	res := seq.ForEach(func(s S) {
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

	if lenOk && sz != i {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, i)
	}

	if err := Error(res); err != nil {
		t.Errorf("Seq returned error: %s", err)
	}
	if sz, lenOk = res.Len(); !lenOk || sz != 0 {
		t.Errorf("Seq returned non-empty from ForEach")
	}
}

func seqIsForEachIndex[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz, lenOk := seq.Len()
	if lenOk {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	count := 0
	res := seq.ForEachIndex(func(i int, s S) {
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

	if lenOk && sz != count {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, count)
	}

	if err := Error(res); err != nil {
		t.Errorf("Seq returned error: %s", err)
	}
	if sz, lenOk = res.Len(); !lenOk || sz != 0 {
		t.Errorf("Seq returned non-empty from ForEachIndex")
	}
}

func seqIsTake[S comparable](t *testing.T, createSeq func() Seq[S], ss []S) {
	t.Helper()

	seq := createSeq()
	sz, lenOk := seq.Len()
	if lenOk {
		t.Run("Len", func(t *testing.T) {
			if sz != len(ss) {
				t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
			}
		})
	}

	t.Run("0", func(t *testing.T) {
		seq = createSeq()
		head, tail := seq.Take(0)
		if len(head) != 0 {
			t.Errorf("When calling Take(0) 'head' must be the empty array")
		}
		if tailLen, _ := tail.Len(); lenOk && tailLen != sz {
			t.Errorf("When calling Take(0) 'tail' must have the same length")
		}
	})

	// Ensure we can Take(n) for different n, and rebuild the exact ss
	for _, n := range []int{1, 2, 3, 100} {
		t.Run(fmt.Sprintf("%d", n), func(t *testing.T) {
			seq = createSeq()
			count := 0
			for head, tail := seq.Take(n); len(head) != 0; head, tail = tail.Take(n) {
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

func seqIsTakeWhile[S comparable](t *testing.T, createSeq func() Seq[S], ss []S) {
	t.Helper()

	seq := createSeq()
	sz, lenOk := seq.Len()
	if lenOk {
		t.Run("Len", func(t *testing.T) {
			if sz != len(ss) {
				t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
			}
		})
	}

	t.Run("skip-all", func(t *testing.T) {
		seq = createSeq()
		head, tail := seq.TakeWhile(func(s S) bool { return false })
		if len(head) != 0 {
			t.Errorf("When calling TakeWhile(false) 'head' must be the empty array")
		}
		if tailLen, _ := tail.Len(); lenOk && sz != tailLen {
			t.Errorf("When calling TakeWhile(false) 'tail' must have the same length. Expected %d, found %d", sz, tailLen)
		}
	})

	t.Run("all", func(t *testing.T) {
		seq = createSeq()
		head, tail := seq.TakeWhile(func(s S) bool { return true })
		if lenOk && len(head) != sz {
			t.Errorf("When calling TakeWhile(true) 'head' must be the entire array. len(head)=%d", len(head))
		}
		if remaining, _ := tail.First(); remaining.Ok() {
			t.Errorf("When calling TakeWhile(true) 'tail' must be empty. Found %v", remaining.val)
		}
	})
}

func seqSkip[S comparable](t *testing.T, createSeq func() Seq[S], ss []S) {
	t.Helper()

	sz, lenOk := createSeq().Len()

	t.Run("skip-all", func(t *testing.T) {
		seq := createSeq()
		tail := seq.Skip(100_000)
		if tailLen, _ := tail.Len(); lenOk && tailLen != 0 {
			t.Errorf("When calling Skip(100,000) 'tail' must be empty")
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after Skipping everything")
		}
	})

	t.Run("one-at-a-time", func(t *testing.T) {
		i := 0
		seq := createSeq()
		for ; i < 100_000; i++ {
			seq = seq.Skip(1)
			if l, _ := seq.Len(); l == 0 {
				break
			}
		}
		if sz != LenUnknown && i == 100_000 {
			t.Errorf("Failed to Skip() Seq 1-by-1, never became empty")
		}
		if fst, _ := seq.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after Skipping everything")
		}
	})

	if sz == LenUnknown {
		return // rest of these tests require a Len
	}

	t.Run("Len", func(t *testing.T) {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	})

	t.Run("skip-all-exact", func(t *testing.T) {
		seq := createSeq()
		tail := seq.Skip(sz)
		if tailLen, _ := tail.Len(); lenOk && tailLen != 0 {
			t.Errorf("When calling Skip(sz) 'tail' must be empty")
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after Skipping everything")
		}
	})
}

func seqIsWhere[S comparable](t *testing.T, createSeq func() Seq[S], ss []S) {
	t.Helper()

	t.Run("false/first", func(t *testing.T) {
		seq := createSeq()
		wh := seq.Where(func(_ S) bool { return false })
		if fst, _ := wh.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after dropping everything with where=false")
		}
	})

	t.Run("false/array", func(t *testing.T) {
		seq := createSeq()
		wh := seq.Where(func(_ S) bool { return false }).Array()
		if l, _ := wh.Len(); l != 0 {
			t.Errorf("Must create empty array after dropping everything with where=false")
		}
	})

	t.Run("false/take-while", func(t *testing.T) {
		seq := createSeq()
		head, tail := seq.Where(func(_ S) bool { return false }).TakeWhile(func(_ S) bool { return true })
		if len(head) != 0 {
			t.Errorf("Must create empty array after dropping everything with where=false")
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() from tail, after dropping everything with where=false")
		}
	})

	t.Run("true/first", func(t *testing.T) {
		seq := createSeq()
		wh := seq.Where(func(_ S) bool { return true })
		i := 0
		for fst, tail := wh.First(); fst.Ok(); fst, tail = tail.First() {
			if ss[i] != fst.val {
				t.Errorf("Unexpected value at inde %d. Expected %v, got %v", i, ss[i], fst.val)
			}
			i++
		}
		if i != len(ss) {
			t.Errorf("Unexpected number of elements in Seq.Where(true). Expected %d, got %d", len(ss), i)
		}
	})

	t.Run("true/array", func(t *testing.T) {
		seq := createSeq()
		arr := seq.Where(func(_ S) bool { return true }).Array()

		if len(arr) != len(ss) {
			t.Errorf("Unexpected number of elements in Seq.Where(true). Expected %d, got %d", len(ss), len(arr))
		}
		if !reflect.DeepEqual(arr.AsSlice(), ss) {
			t.Errorf("Array elements mismatch. Expected %v, got %v", ss, arr)
		}
	})

	t.Run("true/take-while", func(t *testing.T) {
		seq := createSeq()
		head, tail := seq.Where(func(_ S) bool { return true }).TakeWhile(func(_ S) bool { return true })
		if len(head) != len(ss) {
			t.Errorf("Must create empty array after dropping everything with where=false")
		}
		if !reflect.DeepEqual(head.AsSlice(), ss) {
			t.Errorf("Array elements mismatch. Expected %v, got %v", ss, head)
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Tail should be empty. Got %v", fst.val)
		}
	})
}

func seqIsFirst[S comparable](t *testing.T, seq Seq[S], ss []S) {
	t.Helper()

	sz, lenOk := seq.Len()
	if lenOk {
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

func seqIsAll[S comparable](t *testing.T, createSeq func() Seq[S], ss []S) {
	t.Helper()

	t.Run("false", func(t *testing.T) {
		if All(createSeq(), func(_ S) bool { return false }) {
			t.Errorf("All(false) should be false")
		}
	})

	t.Run("true", func(t *testing.T) {
		if !All(createSeq(), func(_ S) bool { return true }) {
			t.Errorf("All(true) should be true")
		}
	})
}

func seqIsAny[S comparable](t *testing.T, createSeq func() Seq[S], ss []S) {
	t.Helper()

	t.Run("false", func(t *testing.T) {
		if Any(createSeq(), func(_ S) bool { return false }) {
			t.Errorf("Any(false) should be false")
		}
	})

	t.Run("true", func(t *testing.T) {
		if !Any(createSeq(), func(_ S) bool { return true }) {
			t.Errorf("Any(true) should be true")
		}
	})
}
