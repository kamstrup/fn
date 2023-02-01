package seq

import "github.com/kamstrup/fn/opt"

var _ Seq[int] = whileSeq[int]{}

type whileSeq[T any] struct {
	seq  Seq[T]
	pred Predicate[T]
}

func WhileOf[T any](seq Seq[T], pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  seq,
		pred: pred,
	}
}

func (w whileSeq[T]) ForEach(f Func1[T]) Seq[T] {
	var (
		fst  opt.Opt[T]
		tail Seq[T]
	)
	for fst, tail = w.First(); fst.Ok(); fst, tail = tail.First() {
		f(fst.Must())
	}

	return tail
}

func (w whileSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	var (
		fst  opt.Opt[T]
		tail Seq[T]
		i    int
	)
	for fst, tail = w.First(); fst.Ok(); fst, tail = tail.First() {
		f(i, fst.Must())
		i++
	}

	return tail
}

func (w whileSeq[T]) Len() (int, bool) {
	return LenUnknown, false
}

func (w whileSeq[T]) Values() Slice[T] {
	head, _ := w.seq.TakeWhile(w.pred)
	return head
}

func (w whileSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n == 0 {
		return []T{}, w
	}

	var (
		arr  []T
		fst  opt.Opt[T]
		tail Seq[T] = w
	)

	for i := 0; i < n; i++ {
		fst, tail = tail.First()
		val, err := fst.Return()
		if err != nil {
			return arr, ErrorOf[T](err)
		}
		arr = append(arr, val)
	}
	return arr, tail
}

func (w whileSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	var (
		arr  []T
		fst  opt.Opt[T]
		tail Seq[T]
	)
	for fst, tail = w.First(); fst.Ok(); fst, tail = tail.First() {
		val := fst.Must()
		if pred(val) {
			arr = append(arr, val)
		} else {
			return arr, PrependOf(val, tail)
		}
	}

	return arr, tail
}

func (w whileSeq[T]) Skip(n int) Seq[T] {
	var (
		i    int
		fst  opt.Opt[T]
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

func (w whileSeq[T]) First() (opt.Opt[T], Seq[T]) {
	fst, tail := w.seq.First()
	val, err := fst.Return()
	if err != nil {
		return fst, tail
	}
	if !w.pred(val) {
		return opt.Empty[T](), errOrEmpty(tail)
	}
	return fst, whileSeq[T]{
		seq:  tail,
		pred: w.pred,
	}
}

func (w whileSeq[K]) Map(shaper FuncMap[K, K]) Seq[K] {
	return mappedSeq[K, K]{
		f:   shaper,
		seq: w,
	}
}

func (w whileSeq[K]) Error() error {
	return Error(w.seq)
}
