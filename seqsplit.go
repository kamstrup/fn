package fn

import "github.com/kamstrup/fn/opt"

// SplitChoice is the return type for FuncSplit used with SplitOf.
// It determines how FuncSplit splits a Seq into sub-seqs.
type SplitChoice uint8

const (
	// SplitKeep indicates that the current element should be included in the current sub-seq.
	SplitKeep SplitChoice = iota
	// SplitSeparate indicates that a new sub-seq should be started and current element discarded
	SplitSeparate
	// SplitSeparateKeep indicates that the current element is the final element in the current seq,
	// and that a new sub-seq should be started.
	SplitSeparateKeep
)

// FuncSplit decides how to split a Seq up into sub-seqs.
type FuncSplit[T any] func(val T) SplitChoice

type splitSeq[T any] struct {
	seq   Seq[T]
	split FuncSplit[T]
}

// SplitOf splits a Seq into sub-seqs based on a FuncSplit.
// The splitting algorithm implemented is only "semi-lazy" in the following way:
// Each split is read eagerly into an Slice, but the tail is not executed.
// The split seq needs to work in this semi-lazy way in order to guarantee
// that methods from the Seq interface returning a tail indeed returns a valid tail.
func SplitOf[T any](seq Seq[T], splitter FuncSplit[T]) Seq[Seq[T]] {
	return splitSeq[T]{
		seq:   seq,
		split: splitter,
	}
}

func (s splitSeq[T]) ForEach(f Func1[Seq[T]]) Seq[Seq[T]] {
	var (
		fst  opt.Opt[Seq[T]]
		tail Seq[Seq[T]]
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		f(fst.Must())
	}
	return tail
}

func (s splitSeq[T]) ForEachIndex(f Func2[int, Seq[T]]) Seq[Seq[T]] {
	var (
		fst  opt.Opt[Seq[T]]
		tail Seq[Seq[T]]
		i    = 0
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		f(i, fst.Must())
		i++
	}
	return tail
}

func (s splitSeq[T]) Len() (int, bool) {
	if sz, _ := s.seq.Len(); sz == LenInfinite {
		return LenInfinite, false
	}
	return LenUnknown, false
}

func (s splitSeq[T]) Values() Slice[Seq[T]] {
	var (
		fst  opt.Opt[Seq[T]]
		tail Seq[Seq[T]]
		arr  []Seq[T]
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		arr = append(arr, fst.Must())
	}
	return arr
}

func (s splitSeq[T]) Take(n int) (Slice[Seq[T]], Seq[Seq[T]]) {
	if n == 0 {
		return []Seq[T]{}, s
	}

	var (
		arr  []Seq[T]
		fst  opt.Opt[Seq[T]]
		tail Seq[Seq[T]] = s
	)

	for i := 0; i < n; i++ {
		fst, tail = tail.First()
		headVal, headErr := fst.Return()
		if headErr != nil {
			return arr, ErrorOf[Seq[T]](headErr)
		}
		arr = append(arr, headVal)
	}
	return arr, tail
}

func (s splitSeq[T]) TakeWhile(pred Predicate[Seq[T]]) (Slice[Seq[T]], Seq[Seq[T]]) {
	var (
		fst  opt.Opt[Seq[T]]
		tail Seq[Seq[T]]
		arr  []Seq[T]
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		val := fst.Must()
		if pred(val) {
			arr = append(arr, val)
		} else {
			return arr, ConcatOf(SingletOf(val), tail)
		}

	}
	return arr, tail
}

func (s splitSeq[T]) Skip(n int) Seq[Seq[T]] {
	var (
		fst  opt.Opt[Seq[T]]
		tail Seq[Seq[T]]
	)
	for fst, tail = s.First(); fst.Ok() && n > 0; fst, tail = tail.First() {
		n--
	}
	return tail
}

func (s splitSeq[T]) Where(pred Predicate[Seq[T]]) Seq[Seq[T]] {
	return whereSeq[Seq[T]]{
		seq:  s,
		pred: pred,
	}
}

func (s splitSeq[T]) While(pred Predicate[Seq[T]]) Seq[Seq[T]] {
	return whileSeq[Seq[T]]{
		seq:  s,
		pred: pred,
	}
}

func (s splitSeq[T]) First() (opt.Opt[Seq[T]], Seq[Seq[T]]) {
	// TODO: special case for Slice?
	var (
		arr  []T
		fst  opt.Opt[T]
		tail = s.seq
	)
	for {
		fst, tail = tail.First()
		val, err := fst.Return()
		if err != nil {
			if len(arr) > 0 {
				return opt.Of(SliceOf(arr)), ErrorOf[Seq[T]](err)
			} else {
				return opt.Empty[Seq[T]](), ErrorOf[Seq[T]](err)
			}

		}
		switch s.split(val) {
		case SplitKeep:
			arr = append(arr, val)
		case SplitSeparate:
			return opt.Of(SliceOf(arr)), splitSeq[T]{seq: tail, split: s.split}
		case SplitSeparateKeep:
			arr = append(arr, val)
			return opt.Of(SliceOf(arr)), splitSeq[T]{seq: tail, split: s.split}
		default:
			panic("fn: invalid SplitChoice")
		}
	}
}

func (s splitSeq[T]) Map(m FuncMap[Seq[T], Seq[T]]) Seq[Seq[T]] {
	return mappedSeq[Seq[T], Seq[T]]{
		f:   m,
		seq: s,
	}
}

func (s splitSeq[T]) Error() error {
	return Error(s.seq)
}
