package fn

type Opt[T comparable] struct {
	val T
	err error
}

func Try[T comparable](t T, err error) Opt[T] {
	return Opt[T]{t, err}
}

func TryDo[T comparable](f FuncSourceErr[T]) Opt[T] {
	t, err := f()
	return Opt[T]{t, err}
}

func TryDoRecover[T comparable](f FuncSourceErr[T]) (opt Opt[T]) {
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

func (o Opt[T]) OnErr(t T) T {
	if o.err != nil {
		return t
	}
	return o.val
}

func (o Opt[T]) OnErrFn(errFn func(err error) T) T {
	if o.err != nil {
		return errFn(o.err)
	}
	return o.val
}

func (o Opt[T]) Return() (T, error) {
	return o.val, o.err
}
