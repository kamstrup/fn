package fn

type errorSeq[T any] struct {
	error
}

// Error returns an error, if there is an error associated with a Seq.
// A sequence has an associated error if it has a method Error() that returns an error.
func Error[T any](seq Seq[T]) error {
	if errSeq, ok := seq.(interface {
		Error() error
	}); ok {
		return errSeq.Error()
	}

	return nil
}

// ErrorOf creates an "error sequence".
// An error sequence is empty and always returns itself. Calling Seq.First() returns an error Opt.
// Calling Error() on an error sequence returns the wrapped error.
func ErrorOf[T any](err error) Seq[T] {
	return errorSeq[T]{err}
}

// Error returns the wrapped error and implements the behavior defined by the package function Error().
func (e errorSeq[T]) Error() error {
	return e.error
}

func (e errorSeq[T]) ForEach(f Func1[T]) Seq[T] {
	return e
}

func (e errorSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	return e
}

func (e errorSeq[T]) Len() (int, bool) {
	return 0, true
}

func (e errorSeq[T]) Array() Array[T] {
	return nil
}

func (e errorSeq[T]) Take(n int) (Array[T], Seq[T]) {
	return nil, e
}

func (e errorSeq[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	return nil, e
}

func (e errorSeq[T]) Skip(n int) Seq[T] {
	return e
}

func (e errorSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return e
}

func (e errorSeq[T]) While(pred Predicate[T]) Seq[T] {
	return e
}

func (e errorSeq[T]) First() (Opt[T], Seq[T]) {
	return OptErr[T](e.error), e
}
