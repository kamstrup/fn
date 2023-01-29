package fn

import "errors"

var ErrEmpty = errors.New("empty")

// Opt is a light wrapper around a value or an error.
// Opts should always be passed by value. If you see *Opt anywhere something is wrong.
type Opt[T any] struct {
	val T
	err error
}

// OptMap converts an option into some other type.
// If you want to keep the same type it may be easier to use Opt.Map.
func OptMap[S, T any](opt Opt[S], f FuncMap[S, T]) Opt[T] {
	if opt.err != nil {
		return OptErr[T](opt.err)
	}
	return OptOf(f(opt.val))
}

// OptOf creates a new opt wrapping a value.
func OptOf[T any](t T) Opt[T] {
	return Opt[T]{val: t}
}

// OptErr creates a new opt wrapping an error.
func OptErr[T any](err error) Opt[T] {
	return Opt[T]{err: err}
}

// OptEmpty creates a new empty opt.
// An empty opt stores the special error ErrEmpty and will respond true to Opt.Empty() and false to Opt.Ok.
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

// Must returns the value wrapped by this opt or panics if there isn't one.
func (o Opt[T]) Must() T {
	if o.err != nil {
		panic(o.err)
	}

	return o.val
}

// Error returns the error, if any, held by this opt.
func (o Opt[T]) Error() error {
	return o.err
}

// OnErr calls a function if the opt is an error or returns the value directly.
// The function returns a default value that will be returned from OnErr.
func (o Opt[T]) OnErr(errFn func(err error) T) T {
	if o.err != nil && o.err != ErrEmpty {
		return errFn(o.err)
	}
	return o.val
}

// Return unpacks the opt into a standard (value, error) pair.
func (o Opt[T]) Return() (T, error) {
	return o.val, o.err
}

func (o Opt[T]) Empty() bool {
	return o.err != nil
}

func (o Opt[T]) Ok() bool {
	return o.err == nil
}

// Or returns the value held by this opt or an alternate value if it is empty or an error.
func (o Opt[T]) Or(altValue T) T {
	if o.Empty() {
		return altValue
	}
	return o.val
}

// Seq interprets this opt as a single-valued Seq.
func (o Opt[T]) Seq() Seq[T] {
	if o.err != nil {
		return ErrorOf[T](o.err)
	} else {
		return SingletOf(o.val)
	}
}
