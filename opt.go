package fn

import "errors"

var ErrEmpty = errors.New("empty")

type Opt[T any] struct {
	val T
	err error
}

func OptMap[S, T any](opt Opt[S], f FuncMap[S, T]) Opt[T] {
	if opt.err != nil {
		return OptErr[T](opt.err)
	}
	return OptOf(f(opt.val))
}

func OptOf[T any](t T) Opt[T] {
	return Opt[T]{val: t}
}

func OptErr[T any](err error) Opt[T] {
	return Opt[T]{err: err}
}

func OptEmpty[T any]() Opt[T] {
	return Opt[T]{err: ErrEmpty}
}

// Map applies a function to the value of the Opt, unless the Opt is empty.
// If you need to change the type inside the Opt you will have to use OptMap.
func (o Opt[T]) Map(f FuncMap[T, T]) Opt[T] {
	if o.Empty() {
		return o
	}
	return OptOf(f(o.val))
}

func (o Opt[T]) Must() T {
	if o.err != nil {
		panic(o.err)
	}

	return o.val
}

func (o Opt[T]) Error() error {
	return o.err
}

func (o Opt[T]) OnErr(errFn func(err error) T) T {
	if o.err != nil && o.err != ErrEmpty {
		return errFn(o.err)
	}
	return o.val
}

func (o Opt[T]) Return() (T, error) {
	return o.val, o.err
}

func (o Opt[T]) Empty() bool {
	return o.err != nil
}

func (o Opt[T]) Ok() bool {
	return o.err == nil
}

func (o Opt[T]) Or(altValue T) T {
	if o.Empty() {
		return altValue
	}
	return o.val
}

func (o Opt[T]) Seq() Seq[T] {
	if o.err != nil {
		return ErrorOf[T](o.err)
	} else {
		return SingletOf(o.val)
	}
}
