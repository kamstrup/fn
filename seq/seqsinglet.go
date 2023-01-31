package seq

import (
	"github.com/kamstrup/fn/opt"
)

type singletSeq[T any] struct {
	val T
}

func SingletOf[T any](t T) Seq[T] {
	return singletSeq[T]{t}
}

func (s singletSeq[T]) ForEach(f Func1[T]) Seq[T] {
	f(s.val)
	return SeqEmpty[T]()
}

func (s singletSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	f(0, s.val)
	return SeqEmpty[T]()
}

func (s singletSeq[T]) Len() (int, bool) {
	return 1, true
}

func (s singletSeq[T]) Values() Slice[T] {
	return Slice[T]{s.val}
}

func (s singletSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n == 0 {
		return Slice[T]{}, s
	}
	return Slice[T]{s.val}, SeqEmpty[T]()
}

func (s singletSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	if pred(s.val) {
		return Slice[T]{s.val}, SeqEmpty[T]()
	}
	return Slice[T]{}, s
}

func (s singletSeq[T]) Skip(n int) Seq[T] {
	if n == 0 {
		return s
	}
	return SeqEmpty[T]()
}

func (s singletSeq[T]) Where(pred Predicate[T]) Seq[T] {
	if pred(s.val) { // eager impl, for optimization
		return s
	}
	return SeqEmpty[T]()
}

func (s singletSeq[T]) While(pred Predicate[T]) Seq[T] {
	if pred(s.val) { // eager impl, for optimization
		return s
	}
	return SeqEmpty[T]()
}

func (s singletSeq[T]) First() (opt.Opt[T], Seq[T]) {
	return opt.Of(s.val), SeqEmpty[T]()
}

func (s singletSeq[T]) Map(f FuncMap[T, T]) Seq[T] {
	return MappingOf[T, T](s, f)
}

func (s singletSeq[T]) Error() error {
	return Error(s.val)
}
