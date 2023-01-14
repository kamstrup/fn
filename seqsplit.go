package fn

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
// Each split is read eagerly into an Array, but the tail is not executed.
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
		fst  Opt[Seq[T]]
		tail Seq[Seq[T]]
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		f(fst.val)
	}
	return tail
}

func (s splitSeq[T]) ForEachIndex(f Func2[int, Seq[T]]) Seq[Seq[T]] {
	var (
		fst  Opt[Seq[T]]
		tail Seq[Seq[T]]
		i    = 0
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		f(i, fst.val)
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

func (s splitSeq[T]) Array() Array[Seq[T]] {
	var (
		fst  Opt[Seq[T]]
		tail Seq[Seq[T]]
		arr  []Seq[T]
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		arr = append(arr, fst.val)
	}
	return arr
}

func (s splitSeq[T]) Take(n int) (Array[Seq[T]], Seq[Seq[T]]) {
	if n == 0 {
		return []Seq[T]{}, s
	}

	var (
		arr  []Seq[T]
		fst  Opt[Seq[T]]
		tail Seq[Seq[T]] = s
	)

	for i := 0; i < n; i++ {
		fst, tail = tail.First()
		if fst.Empty() {
			return arr, tail
		}
		arr = append(arr, fst.val)
	}
	return arr, tail
}

func (s splitSeq[T]) TakeWhile(pred Predicate[Seq[T]]) (Array[Seq[T]], Seq[Seq[T]]) {
	var (
		fst  Opt[Seq[T]]
		tail Seq[Seq[T]]
		arr  []Seq[T]
	)
	for fst, tail = s.First(); fst.Ok(); fst, tail = tail.First() {
		if pred(fst.val) {
			arr = append(arr, fst.val)
		} else {
			return arr, ConcatOf(SingletOf(fst.val), tail)
		}

	}
	return arr, tail
}

func (s splitSeq[T]) Skip(n int) Seq[Seq[T]] {
	var (
		fst  Opt[Seq[T]]
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

func (s splitSeq[T]) First() (Opt[Seq[T]], Seq[Seq[T]]) {
	// TODO: special case for Array?
	var (
		arr  []T
		fst  Opt[T]
		tail = s.seq
	)
	for {
		fst, tail = tail.First()
		if fst.Empty() {
			if len(arr) > 0 {
				return OptOf(ArrayOf(arr)), SeqEmpty[Seq[T]]()
			} else {
				return OptEmpty[Seq[T]](), SeqEmpty[Seq[T]]()
			}

		}
		switch s.split(fst.val) {
		case SplitKeep:
			arr = append(arr, fst.val)
		case SplitSeparate:
			return OptOf(ArrayOf(arr)), splitSeq[T]{seq: tail, split: s.split}
		case SplitSeparateKeep:
			arr = append(arr, fst.val)
			return OptOf(ArrayOf(arr)), splitSeq[T]{seq: tail, split: s.split}
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
