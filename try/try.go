package fntry

import (
	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/opt"
)

type FuncSourceErr[T any] func() (T, error)

// Of creates a new fn.Opt from a value and an error
func Of[T any](t T, err error) opt.Opt[T] {
	if err != nil {
		return opt.ErrorOf[T](err)
	}
	return opt.Of(t)
}

// Call a function returning a fn.Opt with the result
func Call[T any](f FuncSourceErr[T]) opt.Opt[T] {
	t, err := f()
	return Of(t, err)
}

// CallRecover a function returning a fn.Opt with the result.
// If the function panics it is recovered and returned as ErrPanic.
func CallRecover[T any](f FuncSourceErr[T]) (op opt.Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			op = opt.ErrorOf[T](ErrPanic{V: r})
		}
	}()

	return Call(f)
}

// Apply calls a function with an argument and returns the result wrapped in a fn.Opt.
func Apply[S any, T any](f fn.FuncMapErr[S, T], s S) opt.Opt[T] {
	t, err := f(s)
	return Of(t, err)
}

// ApplyRecover calls a function with an argument and returns the result wrapped in a fn.Opt.
// If the function panics it is recovered and returned as ErrPanic.
func ApplyRecover[S any, T any](f fn.FuncMapErr[S, T], s S) (op opt.Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			op = opt.ErrorOf[T](ErrPanic{V: r})
		}
	}()

	return Apply(f, s)
}
