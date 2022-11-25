package fn

var _ Seq[int] = whileSeq[int]{}

type whileSeq[T any] struct {
	seq  Seq[T]
	pred Predicate[T]
}

func (w whileSeq[T]) ForEach(f Func1[T]) {
	for fst, tail := w.First(); !fst.Empty(); fst, tail = tail.First() {
		f(fst.val)
	}
}

func (w whileSeq[T]) ForEachIndex(f Func2[int, T]) {
	i := 0
	for fst, tail := w.First(); !fst.Empty(); fst, tail = tail.First() {
		f(i, fst.val)
		i++
	}
}

func (w whileSeq[T]) Len() (int, bool) {
	return LenUnknown, false
}

func (w whileSeq[T]) Array() Array[T] {
	head, _ := w.seq.TakeWhile(w.pred)
	return head
}

func (w whileSeq[T]) Take(n int) (Array[T], Seq[T]) {
	if n == 0 {
		return []T{}, w
	}

	i := 0
	var arr []T
	for fst, tail := w.First(); !fst.Empty() && i < n; fst, tail = tail.First() {
		arr = append(arr, fst.val)
		i++
	}
	return arr, SeqEmpty[T]()
}

func (w whileSeq[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	var (
		arr  []T
		fst  Opt[T]
		tail Seq[T]
	)
	for fst, tail = w.First(); !fst.Empty() && pred(fst.val); fst, tail = tail.First() {
		arr = append(arr, fst.val)
	}
	return arr, tail
}

func (w whileSeq[T]) Skip(n int) Seq[T] {
	var (
		i    int
		fst  Opt[T]
		tail Seq[T]
	)
	for fst, tail = w.First(); !fst.Empty() && i < n; fst, tail = tail.First() {
		i++
	}
	return tail
}

func (w whileSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  w,
		pred: pred,
	}
}

func (w whileSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  w,
		pred: pred,
	}
}

func (w whileSeq[T]) First() (Opt[T], Seq[T]) {
	fst, tail := w.seq.First()
	if fst.Empty() || !w.pred(fst.val) {
		return OptEmpty[T](), SeqEmpty[T]()
	}
	return fst, whileSeq[T]{
		seq:  tail,
		pred: w.pred,
	}
}

func (w whileSeq[T]) All(pred Predicate[T]) bool {
	return seqAll[T](w, pred)
}

func (w whileSeq[T]) Any(pred Predicate[T]) bool {
	return seqAny[T](w, pred)
}
