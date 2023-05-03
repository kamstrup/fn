package seq

import "github.com/kamstrup/fn/opt"

type errorSeq[T any] struct {
	error
}

// ErrorOf creates an "error sequence".
// An error sequence is empty and always returns itself. Calling Seq.First() returns an error Opt.
// If the error argument is nil or opt.ErrEmpty a non-error empty seq is returned.
func ErrorOf[T any](err error) Seq[T] {
	if err == opt.ErrEmpty || err == nil {
		return Empty[T]()
	}
	return errorSeq[T]{err}
}

// Error returns the wrapped error and implements the behavior defined by the package function Error().
func (e errorSeq[T]) Error() error {
	return e.error
}

func (e errorSeq[T]) ForEach(f Func1[T]) opt.Opt[T] {
	return opt.ErrorOf[T](e.error)
}

func (e errorSeq[T]) ForEachIndex(f Func2[int, T]) opt.Opt[T] {
	return opt.ErrorOf[T](e.error)
}

func (e errorSeq[T]) Len() (int, bool) {
	return 0, true
}

func (e errorSeq[T]) ToSlice() Slice[T] {
	return nil
}

func (e errorSeq[T]) Limit(n int) Seq[T] {
	return e
}

func (e errorSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n < 0 {
		panic("must take >= 0 elements")
	}
	return nil, e
}

func (e errorSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	return nil, e
}

func (e errorSeq[T]) Skip(n int) Seq[T] {
	if n < 0 {
		panic("must skip >= 0 elements")
	}
	return e
}

func (e errorSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return e
}

func (e errorSeq[T]) While(pred Predicate[T]) Seq[T] {
	return e
}

func (e errorSeq[T]) First() (opt.Opt[T], Seq[T]) {
	return opt.ErrorOf[T](e.error), e
}

func (e errorSeq[T]) Map(shaper FuncMap[T, T]) Seq[T] {
	return e
}
