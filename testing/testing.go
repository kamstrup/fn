package fntesting

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

type TestSeq[S comparable] struct {
	t  *testing.T
	sq seq.Seq[S]
}

type Suite[S any] struct {
	t         *testing.T
	createSeq func() seq.Seq[S]
	equal     func(S, S) bool
}

type TestOpt[S comparable] struct {
	t   *testing.T
	opt opt.Opt[S]
}

func TestOf[S comparable](t *testing.T, seq seq.Seq[S]) TestSeq[S] {
	return TestSeq[S]{
		t:  t,
		sq: seq,
	}
}

func OptOf[S comparable](t *testing.T, opt opt.Opt[S]) TestOpt[S] {
	return TestOpt[S]{
		t:   t,
		opt: opt,
	}
}

func SuiteOf[S any](t *testing.T, createSeq func() seq.Seq[S]) Suite[S] {
	return Suite[S]{
		t:         t,
		createSeq: createSeq,
		equal:     func(s1, s2 S) bool { return reflect.DeepEqual(s1, s2) },
	}
}

func (ts TestSeq[S]) LenIs(n int) {
	ts.t.Helper()

	if sz, _ := ts.sq.Len(); sz != n {
		ts.t.Errorf("Seq len mismatch. Expected %d, found %d", n, sz)
	}
}

func (ts TestSeq[S]) Is(ss ...S) {
	ts.t.Helper()

	sz, lenOk := ts.sq.Len()
	if lenOk {
		if sz != len(ss) {
			ts.t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	count := 0
	ts.sq.ForEachIndex(func(i int, s S) {
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
	ts.sq.ForEach(func(s S) {
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

func (to TestOpt[S]) IsError(expectErr any) {
	to.t.Helper()

	val, err := to.opt.Return()
	if !reflect.DeepEqual(err, expectErr) {
		to.t.Errorf("unexpected error: %v (value: %v)", err, val)
	}
}

func (ts Suite[S]) WithEqual(eq func(s1, s2 S) bool) Suite[S] {
	return Suite[S]{
		t:         ts.t,
		createSeq: ts.createSeq,
		equal:     eq,
	}
}

func (ts Suite[S]) Is(ss ...S) {
	ts.t.Helper()

	ts.t.Run("ForEach", func(t *testing.T) {
		ts.seqIsForEach(t, ss)
	})

	ts.t.Run("ForEachIndex", func(t *testing.T) {
		ts.seqIsForEachIndex(t, ss)
	})

	ts.t.Run("Take", func(t *testing.T) {
		ts.seqIsTake(t, ss)
	})

	ts.t.Run("TakeWhile", func(t *testing.T) {
		ts.seqIsTakeWhile(t, ss)
	})

	ts.t.Run("Skip", func(t *testing.T) {
		ts.seqSkip(t, ss)
	})

	ts.t.Run("Where", func(t *testing.T) {
		ts.seqIsWhere(t, ss)
	})
	// TODO: test While(pred)

	ts.t.Run("First", func(t *testing.T) {
		ts.seqIsFirst(t, ss)
	})

	ts.t.Run("All", func(t *testing.T) {
		ts.seqIsAll(t, ss)
	})

	ts.t.Run("Any", func(t *testing.T) {
		ts.seqIsAny(t, ss)
	})
}

func (ts Suite[S]) IsEmpty() {
	ts.t.Helper()

	count := 0
	ts.createSeq().ForEach(func(s S) {
		ts.t.Helper()
		count++
	})
	if count != 0 {
		ts.t.Errorf("Seq is not empty. Length %d", count)
	}

	if arr := ts.createSeq().ToSlice(); len(arr) != 0 {
		ts.t.Errorf("Seq.Slice is not empty. Length %d", len(arr))
	}

	// TODO: more checks for emptiness!
}

func (ts Suite[S]) seqIsForEach(t *testing.T, ss []S) {
	t.Helper()

	sq := ts.createSeq()

	sz, lenOk := sq.Len()
	if lenOk {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	i := 0
	res := sq.ForEach(func(s S) {
		t.Helper()

		if i >= len(ss) {
			t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, i, s)
		} else if !ts.equal(ss[i], s) {
			t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				i, ss[i], s)
		}

		i++
	})

	if lenOk && sz != i {
		t.Errorf("Number of elements in ForEach incorrect. Expected %d, found %d",
			sz, i)
	}

	if err := res.Error(); err != nil {
		t.Errorf("Seq.ForEach returned error: %s", err)
	}
}

func (ts Suite[S]) seqIsForEachIndex(t *testing.T, ss []S) {
	t.Helper()

	sq := ts.createSeq()

	sz, lenOk := sq.Len()
	if lenOk {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	count := 0
	res := sq.ForEachIndex(func(i int, s S) {
		t.Helper()

		count++
		if i >= len(ss) {
			t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, i, s)
		} else if !ts.equal(ss[i], s) {
			t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				i, ss[i], s)
		}
	})

	if lenOk && sz != count {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, count)
	}

	if err := res.Error(); err != nil {
		t.Errorf("Seq.ForEachIndex returned error: %s", err)
	}
}

func (ts Suite[S]) seqIsTake(t *testing.T, ss []S) {
	t.Helper()

	sq := ts.createSeq()
	sz, lenOk := sq.Len()
	if lenOk {
		t.Run("Len", func(t *testing.T) {
			if sz != len(ss) {
				t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
			}
		})
	}

	t.Run("0", func(t *testing.T) {
		sq = ts.createSeq()
		head, tail := sq.Take(0)
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
			sq = ts.createSeq()
			count := 0
			for head, tail := sq.Take(n); len(head) != 0; head, tail = tail.Take(n) {
				for i := range head {
					if i+count >= len(ss) {
						t.Fatalf("Seq index out of bounds. Expected %d items, found %v at index %d",
							len(ss), head[i], i+count)
					}
					if !ts.equal(ss[count+i], head[i]) {
						t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
							count+i, ss[count+i], head[i])
					}
				}

				count += len(head)
				if count > len(ss) {
					t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
						len(ss)-1, count, head[0])
				}
			}

			if sz != seq.LenUnknown && sz != count {
				t.Errorf("Number of elements in Take(%d) incorrect. Expected %d, found %d",
					n, sz, count)
			}
		})
	}
}

func (ts Suite[S]) seqIsTakeWhile(t *testing.T, ss []S) {
	t.Helper()

	sq := ts.createSeq()
	sz, lenOk := sq.Len()
	if lenOk {
		t.Run("Len", func(t *testing.T) {
			if sz != len(ss) {
				t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
			}
		})
	}

	t.Run("skip-all", func(t *testing.T) {
		sq = ts.createSeq()
		head, tail := sq.TakeWhile(func(s S) bool { return false })
		if len(head) != 0 {
			t.Errorf("When calling TakeWhile(false) 'head' must be the empty array")
		}
		if tailLen, _ := tail.Len(); lenOk && sz != tailLen {
			t.Errorf("When calling TakeWhile(false) 'tail' must have the same length. Expected %d, found %d", sz, tailLen)
		}
	})

	t.Run("all", func(t *testing.T) {
		sq = ts.createSeq()
		head, tail := sq.TakeWhile(func(s S) bool { return true })
		if lenOk && len(head) != sz {
			t.Errorf("When calling TakeWhile(true) 'head' must be the entire array. len(head)=%d", len(head))
		}
		if remaining, _ := tail.First(); remaining.Ok() {
			t.Errorf("When calling TakeWhile(true) 'tail' must be empty. Found %v", remaining.Must())
		}
	})

	t.Run("one-by-one", func(t *testing.T) {
		tail := ts.createSeq()
		var head seq.Slice[S]
		for i := range ss {
			numTakes := 0
			head, tail = tail.TakeWhile(func(s S) bool {
				numTakes++
				return numTakes == 1
			})
			if len(head) != 1 {
				t.Errorf("Should see exactly 1 element, found: %d", len(head))
			}
			if !reflect.DeepEqual(head[0], ss[i]) {
				t.Errorf("When calling TakeWhile(), found unexpected element at index %d:\nFound:    %v\nExpected: %v", i, head[0], ss[i])
			}
		}

		// tail should now be empty
		if remaining, _ := tail.First(); remaining.Ok() {
			t.Errorf("When calling TakeWhile(true) 'tail' must be empty. Found %v", remaining.Must())
		}
	})
}

func (ts Suite[S]) seqSkip(t *testing.T, ss []S) {
	t.Helper()

	sz, lenOk := ts.createSeq().Len()

	t.Run("skip-all", func(t *testing.T) {
		sq := ts.createSeq()
		tail := sq.Skip(100_000)
		if tailLen, _ := tail.Len(); lenOk && tailLen != 0 {
			t.Errorf("When calling Skip(100,000) 'tail' must be empty")
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after Skipping everything")
		}
	})

	t.Run("one-at-a-time", func(t *testing.T) {
		i := 0
		sq := ts.createSeq()
		for ; i < 100_000; i++ {
			sq = sq.Skip(1)
			if l, _ := sq.Len(); l == 0 {
				break
			}
		}
		if sz != seq.LenUnknown && i == 100_000 {
			t.Errorf("Failed to Skip() Seq 1-by-1, never became empty")
		}
		if fst, _ := sq.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after Skipping everything")
		}
	})

	if sz == seq.LenUnknown {
		return // rest of these tests require a Len
	}

	t.Run("Len", func(t *testing.T) {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	})

	t.Run("skip-all-exact", func(t *testing.T) {
		seq := ts.createSeq()
		tail := seq.Skip(sz)
		if tailLen, _ := tail.Len(); lenOk && tailLen != 0 {
			t.Errorf("When calling Skip(sz) 'tail' must be empty")
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after Skipping everything")
		}
	})
}

func (ts Suite[S]) seqIsWhere(t *testing.T, ss []S) {
	t.Helper()

	t.Run("false/first", func(t *testing.T) {
		seq := ts.createSeq()
		wh := seq.Where(func(_ S) bool { return false })
		if fst, _ := wh.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() after dropping everything with where=false")
		}
	})

	t.Run("false/array", func(t *testing.T) {
		seq := ts.createSeq()
		wh := seq.Where(func(_ S) bool { return false }).ToSlice()
		if l, _ := wh.Len(); l != 0 {
			t.Errorf("Must create empty array after dropping everything with where=false")
		}
	})

	t.Run("false/take-while", func(t *testing.T) {
		seq := ts.createSeq()
		head, tail := seq.Where(func(_ S) bool { return false }).TakeWhile(func(_ S) bool { return true })
		if len(head) != 0 {
			t.Errorf("Must create empty array after dropping everything with where=false")
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Must not be able to take First() from tail, after dropping everything with where=false")
		}
	})

	t.Run("true/first", func(t *testing.T) {
		seq := ts.createSeq()
		wh := seq.Where(func(_ S) bool { return true })
		i := 0
		for fst, tail := wh.First(); fst.Ok(); fst, tail = tail.First() {
			val := fst.Must()
			if !ts.equal(ss[i], val) {
				t.Errorf("Unexpected value at inde %d. Expected %v, got %v", i, ss[i], val)
			}
			i++
		}
		if i != len(ss) {
			t.Errorf("Unexpected number of elements in Seq.Where(true). Expected %d, got %d", len(ss), i)
		}
	})

	t.Run("true/array", func(t *testing.T) {
		seq := ts.createSeq()
		arr := seq.Where(func(_ S) bool { return true }).ToSlice()

		if len(arr) != len(ss) {
			t.Errorf("Unexpected number of elements in Seq.Where(true). Expected %d, got %d", len(ss), len(arr))
		}
		if !reflect.DeepEqual([]S(arr), ss) {
			t.Errorf("Slice elements mismatch.\nExpected %v,\ngot      %v", ss, arr)
		}
	})

	t.Run("true/take-while", func(t *testing.T) {
		seq := ts.createSeq()
		head, tail := seq.Where(func(_ S) bool { return true }).TakeWhile(func(_ S) bool { return true })
		if len(head) != len(ss) {
			t.Errorf("Must create empty array after dropping everything with where=false")
		}
		if !reflect.DeepEqual([]S(head), ss) {
			t.Errorf("Slice elements mismatch.\nExpected: %v\nGot     : %v", ss, head)
		}
		if fst, _ := tail.First(); fst.Ok() {
			t.Errorf("Tail should be empty. Got %v", fst.Must())
		}
	})
}

func (ts Suite[S]) seqIsFirst(t *testing.T, ss []S) {
	t.Helper()

	sq := ts.createSeq()
	sz, lenOk := sq.Len()
	if lenOk {
		if sz != len(ss) {
			t.Errorf("Seq len mismatch. Expected %d, found %d", len(ss), sz)
		}
	}

	var fst opt.Opt[S]
	idx := 0
	for fst, sq = sq.First(); fst.Ok(); fst, sq = sq.First() {
		s := fst.Must()
		if idx >= len(ss) {
			t.Fatalf("Seq element index out of bounds. Expected max index %d, but got index %d with value %v",
				len(ss)-1, idx, s)
		} else if !ts.equal(ss[idx], s) {
			t.Errorf("Seq element mismatch at index %d. Expected %v, found %v",
				idx, ss[idx], s)
		}
		idx++
	}

	if sz != seq.LenUnknown && sz != idx {
		t.Errorf("Number of elements in ForEachIndex incorrect. Expected %d, found %d",
			sz, idx)
	}
}

func (ts Suite[S]) seqIsAll(t *testing.T, ss []S) {
	t.Helper()

	t.Run("false", func(t *testing.T) {
		if seq.All(ts.createSeq(), func(_ S) bool { return false }) {
			t.Errorf("All(false) should be false")
		}
	})

	t.Run("true", func(t *testing.T) {
		if !seq.All(ts.createSeq(), func(_ S) bool { return true }) {
			t.Errorf("All(true) should be true")
		}
	})
}

func (ts Suite[S]) seqIsAny(t *testing.T, ss []S) {
	t.Helper()

	t.Run("false", func(t *testing.T) {
		if seq.Any(ts.createSeq(), func(_ S) bool { return false }) {
			t.Errorf("Any(false) should be false")
		}
	})

	t.Run("true", func(t *testing.T) {
		if !seq.Any(ts.createSeq(), func(_ S) bool { return true }) {
			t.Errorf("Any(true) should be true")
		}
	})
}
