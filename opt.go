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

func TryOf[T any](t T, err error) Opt[T] {
	return Opt[T]{t, err}
}

func Try[T any](f FuncSourceErr[T]) Opt[T] {
	t, err := f()
	return Opt[T]{t, err}
}

func TryRecover[T any](f FuncSourceErr[T]) (opt Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			opt.err = ErrPanic{V: r}
		}
	}()

	opt.val, opt.err = f()
	return
}

func TryMap[S any, T any](f FuncMapErr[S, T], s S) Opt[T] {
	t, err := f(s)
	return Opt[T]{t, err}
}

func TryMapRecover[S any, T any](f FuncMapErr[S, T], s S) (opt Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			opt.err = ErrPanic{V: r}
		}
	}()

	opt.val, opt.err = f(s)
	return
}

func (o Opt[T]) Must() T {
	if o.err != nil {
		panic(o.err)
	}

	return o.val
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
