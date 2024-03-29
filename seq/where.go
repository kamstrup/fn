package seq

import "github.com/kamstrup/fn/opt"

var _ Seq[int] = whereSeq[int]{}

type whereSeq[T any] struct {
	seq  Seq[T]
	pred Predicate[T]
}

func WhereOf[T any](seq Seq[T], pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  seq,
		pred: pred,
	}
}

func (ws whereSeq[T]) ForEach(f Func1[T]) opt.Opt[T] {
	return ws.seq.ForEach(func(t T) {
		if ws.pred(t) {
			f(t)
		}
	})
}

func (ws whereSeq[T]) ForEachIndex(f Func2[int, T]) opt.Opt[T] {
	i := 0
	return ws.seq.ForEachIndex(func(_ int, t T) {
		if ws.pred(t) {
			f(i, t)
			i++
		}
	})
}

func (ws whereSeq[T]) Len() (int, bool) {
	if sz, _ := ws.seq.Len(); sz == 0 {
		return 0, true
	}
	return LenUnknown, false
}

func (ws whereSeq[T]) ToSlice() Slice[T] {
	if l, _ := ws.Len(); l == 0 {
		return Slice[T](nil)
	}

	var arr []T
	ws.ForEach(func(t T) {
		arr = append(arr, t)
	})

	return arr
}

func (ws whereSeq[T]) Limit(n int) Seq[T] {
	return LimitOf[T](ws, n)
}

func (ws whereSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if l, _ := ws.Len(); l == 0 || n == 0 {
		return Slice[T](nil), Empty[T]()
	}

	// FIXME: TakeWhile + append() does double alloc! Should just be While() + Do()
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

	return arr, tail
}

func (ws whereSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	if l, _ := ws.Len(); l == 0 {
		return Slice[T](nil), Empty[T]()
	}

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

	return arr, tail
}

func (ws whereSeq[T]) Skip(n int) Seq[T] {
	if n == 0 {
		return ws
	}

	var (
		fst  opt.Opt[T]
		tail Seq[T]
		i    = 0
	)
	for fst, tail = ws.First(); !fst.Empty() && i < n; fst, tail = tail.First() {
		i++
	}
	return tail
}

func (ws whereSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  ws,
		pred: pred,
	}
}

func (ws whereSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  ws,
		pred: pred,
	}
}

func (ws whereSeq[T]) First() (opt.Opt[T], Seq[T]) {
	var (
		fst  opt.Opt[T]
		tail Seq[T]
	)
	for fst, tail = ws.seq.First(); ; fst, tail = tail.First() {
		// seek until we find a First element that is true for ws.pred()
		val, err := fst.Return()
		if err != nil {
			return fst, ErrorOf[T](err)
		}
		if ws.pred(val) {
			return fst, whereSeq[T]{tail, ws.pred}
		}
	}
}

func (ws whereSeq[K]) Map(shaper FuncMap[K, K]) Seq[K] {
	return mappedSeq[K, K]{
		f:   shaper,
		seq: ws,
	}
}
