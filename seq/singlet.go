package seq

import (
	"github.com/kamstrup/fn/opt"
)

type singletSeq[T any] struct {
	val T
}

// OptOf returns a Seq from an opt.Opt.
// If the opt is empty the seq is empty, if there is some other error an error seq is returned,
// and if the opt is valid a single element seq is returned.
func OptOf[T any](op opt.Opt[T]) Seq[T] {
	if err := op.Error(); err != nil {
		if err == opt.ErrEmpty {
			return Empty[T]()
		}
		return ErrorOf[T](err)
	}
	return SingletOf(op.Must())
}

// SingletOf returns a single element seq.
func SingletOf[T any](t T) Seq[T] {
	return singletSeq[T]{t}
}

func (s singletSeq[T]) ForEach(f Func1[T]) opt.Opt[T] {
	f(s.val)
	return opt.Zero[T]()
}

func (s singletSeq[T]) ForEachIndex(f Func2[int, T]) opt.Opt[T] {
	f(0, s.val)
	return opt.Zero[T]()
}

func (s singletSeq[T]) Len() (int, bool) {
	return 1, true
}

func (s singletSeq[T]) ToSlice() Slice[T] {
	return Slice[T]{s.val}
}

func (s singletSeq[T]) Limit(n int) Seq[T] {
	return LimitOf[T](s, n)
}

func (s singletSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n == 0 {
		return Slice[T]{}, s
	}
	return Slice[T]{s.val}, Empty[T]()
}

func (s singletSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	if pred(s.val) {
		return Slice[T]{s.val}, Empty[T]()
	}
	return Slice[T]{}, s
}

func (s singletSeq[T]) Skip(n int) Seq[T] {
	if n == 0 {
		return s
	}
	return Empty[T]()
}

func (s singletSeq[T]) Where(pred Predicate[T]) Seq[T] {
	if pred(s.val) { // eager impl, for optimization
		return s
	}
	return Empty[T]()
}

func (s singletSeq[T]) While(pred Predicate[T]) Seq[T] {
	if pred(s.val) { // eager impl, for optimization
		return s
	}
	return Empty[T]()
}

func (s singletSeq[T]) First() (opt.Opt[T], Seq[T]) {
	return opt.Of(s.val), Empty[T]()
}

func (s singletSeq[T]) Map(f FuncMap[T, T]) Seq[T] {
	return MappingOf[T, T](s, f)
}
