package seq

import (
	"github.com/kamstrup/fn/opt"
)

type limitSeq[T any] struct {
	seq   Seq[T]
	limit int
}

func LimitOf[T any](seq Seq[T], limit int) Seq[T] {
	if limit <= 0 {
		return Empty[T]()
	} else if sz, ok := seq.Len(); ok && sz <= limit {
		return seq
	}

	return limitSeq[T]{
		seq:   seq,
		limit: limit,
	}
}

func (l limitSeq[T]) ForEach(f Func1[T]) Seq[T] {
	var (
		i    = 0
		head opt.Opt[T]
		tail Seq[T]
	)

	for head, tail = l.seq.First(); i < l.limit && head.Ok(); head, tail = tail.First() {
		v, _ := head.Return()
		f(v)
		i++
	}

	return tail
}

func (l limitSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	var (
		i    = 0
		head opt.Opt[T]
		tail Seq[T]
	)

	for head, tail = l.seq.First(); i < l.limit && head.Ok(); head, tail = tail.First() {
		v, _ := head.Return()
		f(i, v)
		i++
	}

	return tail
}

func (l limitSeq[T]) Len() (int, bool) {
	sz, ok := l.seq.Len()
	if ok {
		if sz < l.limit {
			return sz, true
		}
		return l.limit, true
	}

	// len not well-defined
	if sz == LenInfinite {
		return l.limit, true
	}
	return LenUnknown, false
}

func (l limitSeq[T]) ToSlice() Slice[T] {
	arr, _ := l.Take(l.limit)
	return arr
}

func (l limitSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n > l.limit {
		n = l.limit
	}

	head, tail := l.seq.Take(n)
	return head, LimitOf(tail, l.limit-len(head))
}

func (l limitSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	i := -1
	head, tail := l.seq.TakeWhile(func(t T) bool {
		i++
		return i < l.limit && pred(t)
	})

	return head, LimitOf(tail, l.limit-i)
}

func (l limitSeq[T]) Skip(n int) Seq[T] {
	if n > l.limit {
		return Empty[T]()
	}
	return LimitOf(l.seq.Skip(n), l.limit-n)
}

func (l limitSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return WhereOf[T](l, pred)
}

func (l limitSeq[T]) While(pred Predicate[T]) Seq[T] {
	return WhileOf[T](l, pred)
}

func (l limitSeq[T]) First() (opt.Opt[T], Seq[T]) {
	fst, tail := l.seq.First()
	return fst, LimitOf(tail, l.limit-1)
}

func (l limitSeq[T]) Map(f FuncMap[T, T]) Seq[T] {
	return MappingOf[T, T](l, f)
}
