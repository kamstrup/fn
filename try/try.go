package fntry

import "github.com/kamstrup/fn"

type FuncSourceErr[T any] func() (T, error)

// Of creates a new fn.Opt from a value and an error
func Of[T any](t T, err error) fn.Opt[T] {
	if err != nil {
		return fn.OptErr[T](err)
	}
	return fn.OptOf(t)
}

// Call a function returning a fn.Opt with the result
func Call[T any](f FuncSourceErr[T]) fn.Opt[T] {
	t, err := f()
	return Of(t, err)
}

// CallRecover a function returning a fn.Opt with the result.
// If the function panics it is recovered and returned as ErrPanic.
func CallRecover[T any](f FuncSourceErr[T]) (opt fn.Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			opt = fn.OptErr[T](ErrPanic{V: r})
		}
	}()

	return Call(f)
}

// Apply calls a function with an argument and returns the result wrapped in a fn.Opt.
func Apply[S any, T any](f fn.FuncMapErr[S, T], s S) fn.Opt[T] {
	t, err := f(s)
	return Of(t, err)
}

// ApplyRecover calls a function with an argument and returns the result wrapped in a fn.Opt.
// If the function panics it is recovered and returned as ErrPanic.
func ApplyRecover[S any, T any](f fn.FuncMapErr[S, T], s S) (opt fn.Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			opt = fn.OptErr[T](ErrPanic{V: r})
		}
	}()

	return Apply(f, s)
}
