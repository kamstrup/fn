package fn

type sourceSeq[T any] struct {
	f FuncSource[T]
}

// SourceOf creates a new Seq yielding the return value of the FuncSource for every element.
// You can for example use it NumbersFrom or Constant.
// Beware when using SourceOf since it often produces a stateful Seq, so order of operations
// may matter.
func SourceOf[T any](f FuncSource[T]) Seq[T] {
	return sourceSeq[T]{f}
}

func (s sourceSeq[T]) ForEach(f Func1[T]) {
	for { // loop infinitely, only a panic or Exit() will stop this
		t := s.f()
		f(t)
	}
}

func (s sourceSeq[T]) ForEachIndex(f Func2[int, T]) {
	for i := 0; ; i++ {
		t := s.f()
		f(i, t)
	}
}

func (s sourceSeq[T]) Len() int {
	return LenUnknown
}

func (s sourceSeq[T]) Array() Array[T] {
	panic("cannot create Array of infinite source")
}

func (s sourceSeq[T]) Take(n int) (Array[T], Seq[T]) {
	arr := make([]T, n)
	for i := 0; i < n; i++ {
		arr[i] = s.f()
	}
	return ArrayOf(arr), s
}

func (s sourceSeq[T]) TakeWhile(predicate Predicate[T]) (Array[T], Seq[T]) {
	var arr []T
	for t := s.f(); predicate(t); t = s.f() {
		arr = append(arr, t)
	}
	return ArrayOf(arr), s
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

func (s sourceSeq[T]) First() (Opt[T], Seq[T]) {
	return OptOf(s.f()), s
}
