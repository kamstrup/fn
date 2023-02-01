package seq

import (
	"github.com/kamstrup/fn/opt"
)

type prepSeq[T any] struct {
	head T
	tail Seq[T]
}

// PrependOf returns a seq that starts with a given element and continues into the provided tail.
// This is not an efficient way of building a large seq of any kind, but mainly intended
// when you need to "unread" or "push back" an element you already executed from a seq.
//
// See also ConcatOf and FlattenOf.
func PrependOf[T any](t T, tail Seq[T]) Seq[T] {
	return prepSeq[T]{
		head: t,
		tail: tail,
	}
}

func (p prepSeq[T]) ForEach(f Func1[T]) Seq[T] {
	f(p.head)
	return p.tail.ForEach(f)
}

func (p prepSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	f(0, p.head)
	i := 1
	return p.tail.ForEach(func(t T) {
		f(i, t)
		i++
	})
}

func (p prepSeq[T]) Len() (int, bool) {
	sz, ok := p.tail.Len()
	if ok {
		return sz + 1, true
	}
	if sz == LenInfinite {
		return LenInfinite, false
	}
	return LenUnknown, false
}

func (p prepSeq[T]) Values() Slice[T] {
	var arr []T
	if sz, ok := p.Len(); ok {
		arr = make([]T, 0, sz)
	}
	p.ForEach(func(t T) {
		arr = append(arr, t)
	})

	return arr
}

func (p prepSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n == 0 {
		return Slice[T]{}, p
	} else if n == 1 {
		return SliceAsArgs(p.head), p.tail
	}

	var arr []T
	if sz, ok := p.Len(); ok {
		if n > sz {
			arr = make([]T, 0, sz)
		} else {
			arr = make([]T, 0, n)
		}
	}

	var (
		i    int
		fst  opt.Opt[T]
		tail Seq[T]
	)
	for fst, tail = p.First(); fst.Ok() && i < n; fst, tail = tail.First() {
		arr = append(arr, fst.Must())
		i++
	}

	return arr, tail
}

func (p prepSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	if !pred(p.head) {
		return Slice[T]{}, p
	}

	var (
		arr  = []T{p.head}
		fst  opt.Opt[T]
		tail Seq[T]
	)
	for fst, tail = p.tail.First(); fst.Ok(); fst, tail = tail.First() {
		val := fst.Must()
		if pred(val) {
			arr = append(arr, val)
		} else {
			return arr, PrependOf(val, tail)
		}
	}

	return arr, tail
}

func (p prepSeq[T]) Skip(n int) Seq[T] {
	if n == 0 {
		return p
	} else if n == 1 {
		return p.tail
	}

	return p.tail.Skip(n - 1)
}

func (p prepSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  p,
		pred: pred,
	}
}

func (p prepSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  p,
		pred: pred,
	}
}

func (p prepSeq[T]) First() (opt.Opt[T], Seq[T]) {
	return opt.Of(p.head), p.tail
}

func (p prepSeq[T]) Map(funcMap FuncMap[T, T]) Seq[T] {
	return mappedSeq[T, T]{
		f:   funcMap,
		seq: p,
	}
}
