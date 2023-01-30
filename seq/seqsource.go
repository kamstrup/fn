package seq

import "github.com/kamstrup/fn/opt"

type sourceSeq[T any] struct {
	f FuncSource[T]
}

// SourceOf creates a new Seq yielding the return value of the FuncSource for every element.
// The length of a source Seq is always LenInfinite.
// Beware when using SourceOf since it often produces a stateful Seq, so order of operations
// may matter.
func SourceOf[T any](f FuncSource[T]) Seq[T] {
	return sourceSeq[T]{f}
}

// Constant returns an infinite Seq that repeats the same value.
func Constant[T any](t T) Seq[T] {
	return SourceOf(func() T {
		return t
	})
}

func (s sourceSeq[T]) ForEach(f Func1[T]) Seq[T] {
	for { // loop infinitely, only a panic or Exit() will stop this
		t := s.f()
		f(t)
	}
}

func (s sourceSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	for i := 0; ; i++ {
		t := s.f()
		f(i, t)
	}
}

func (s sourceSeq[T]) Len() (int, bool) {
	return LenInfinite, false
}

func (s sourceSeq[T]) Values() Slice[T] {
	panic("cannot create Slice of infinite source")
}

func (s sourceSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	arr := make([]T, n)
	for i := 0; i < n; i++ {
		arr[i] = s.f()
	}
	return arr, s
}

func (s sourceSeq[T]) TakeWhile(predicate Predicate[T]) (Slice[T], Seq[T]) {
	var arr []T
	for t := s.f(); predicate(t); t = s.f() {
		arr = append(arr, t)
	}
	return arr, s
}

func (s sourceSeq[T]) Skip(n int) Seq[T] {
	for i := 0; i < n; i++ {
		s.f()
	}
	return s
}

func (s sourceSeq[T]) Where(p Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  s,
		pred: p,
	}
}

func (s sourceSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  s,
		pred: pred,
	}
}

func (s sourceSeq[T]) First() (opt.Opt[T], Seq[T]) {
	return opt.Of(s.f()), s
}

func (s sourceSeq[K]) Map(shaper FuncMap[K, K]) Seq[K] {
	return mappedSeq[K, K]{
		f:   shaper,
		seq: s,
	}
}
