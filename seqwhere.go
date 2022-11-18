package fn

var _ Seq[int] = whereSeq[int]{}

type whereSeq[T comparable] struct {
	seq  Seq[T]
	pred Predicate[T]
}

func (ws whereSeq[T]) ForEach(f Func1[T]) {
	ws.seq.ForEach(func(t T) {
		if ws.pred(t) {
			f(t)
		}
	})
}

func (ws whereSeq[T]) ForEachIndex(f Func2[int, T]) {
	i := 0
	ws.seq.ForEachIndex(func(_ int, t T) {
		if ws.pred(t) {
			f(i, t)
			i++
		}
	})
}

func (ws whereSeq[T]) Len() int {
	if sz := ws.seq.Len(); sz == 0 {
		return 0
	}
	return LenUnknown
}

func (ws whereSeq[T]) Array() Array[T] {
	if ws.Len() == 0 {
		return ArrayOf[T](nil)
	}

	var arr []T
	ws.ForEach(func(t T) {
		arr = append(arr, t)
	})

	return ArrayOf(arr)
}

func (ws whereSeq[T]) Take(n int) (Seq[T], Seq[T]) {
	if ws.Len() == 0 || n == 0 {
		return SeqEmpty[T](), SeqEmpty[T]()
	}

	// TODO: does not really to alloc a slice, if we had a "pulling seq"
	var (
		sz  = 0
		arr []T
	)
	_, tail := ws.TakeWhile(func(t T) bool {
		arr = append(arr, t)
		sz++
		if sz == n {
			return false
		}
		return true
	})

	return ArrayOf[T](arr), tail
}

func (ws whereSeq[T]) TakeWhile(pred Predicate[T]) (Seq[T], Seq[T]) {
	if ws.Len() == 0 {
		return SeqEmpty[T](), SeqEmpty[T]()
	}

	// TODO: does not really to alloc a slice, if we had a "pulling seq"
	var arr []T
	_, tail := ws.seq.TakeWhile(func(t T) bool {
		if !ws.pred(t) {
			return true // skipped by where-clause
		}
		if pred(t) {
			arr = append(arr, t)
			return true
		}
		return false
	})

	return ArrayOf[T](arr), tail
}

func (ws whereSeq[T]) Skip(i int) Seq[T] {
	// TODO implement me
	panic("implement me")
}

func (ws whereSeq[T]) Where(p Predicate[T]) Seq[T] {
	// TODO implement me
	panic("implement me")
}

func (ws whereSeq[T]) First() (Opt[T], Seq[T]) {
	// TODO implement me
	panic("implement me")
}
