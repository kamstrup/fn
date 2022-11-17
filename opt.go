package fn

import "errors"

var errEmpty = errors.New("opt is empty")

type Opt[T comparable] struct {
	val T
	err error
}

func OptMap[S, T comparable](opt Opt[S], f FuncMap[S, T]) Opt[T] {
	if opt.err != nil {
		return OptErr[T](opt.err)
	}
	return OptOf(f(opt.val))
}

func OptOf[T comparable](t T) Opt[T] {
	return Opt[T]{val: t}
}

func OptErr[T comparable](err error) Opt[T] {
	return Opt[T]{err: err}
}

func OptEmpty[T comparable]() Opt[T] {
	return Opt[T]{err: errEmpty}
}

func TryOf[T comparable](t T, err error) Opt[T] {
	return Opt[T]{t, err}
}

func Try[T comparable](f FuncSourceErr[T]) Opt[T] {
	t, err := f()
	return Opt[T]{t, err}
}

func TryRecover[T comparable](f FuncSourceErr[T]) (opt Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			opt.err = ErrPanic{V: r}
		}
	}()

	opt.val, opt.err = f()
	return
}

func TryMap[S any, T comparable](f FuncMapErr[S, T], s S) Opt[T] {
	t, err := f(s)
	return Opt[T]{t, err}
}

func TryMapRecover[S any, T comparable](f FuncMapErr[S, T], s S) (opt Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			opt.err = ErrPanic{V: r}
		}
	}()

	opt.val, opt.err = f(s)
	return
}

func (o Opt[T]) Must() T {
	var defaultVal T
	if o.err != nil {
		panic(o.err)
	} else if o.val == defaultVal {
		panic("unset optional value")
	}

	return o.val
}

func (o Opt[T]) OnErr(errFn func(err error) T) T {
	if o.err != nil && o.err != errEmpty {
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

func (o Opt[T]) Or(altValue T) T {
	if o.Empty() {
		return altValue
	}
	return o.val
}
