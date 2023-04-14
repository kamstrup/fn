package opt

import (
	"errors"
)

// ErrEmpty is a constant error value used to signify when an Opt is empty.
var ErrEmpty = errors.New("empty")

// Opt is a light wrapper around a value or an error.
// Opts should always be passed by value. If you see a pointer to an Opt anywhere something is wrong.
//
// The zero Opt structure holds the default value for T and no error, and is *not* considered empty.
type Opt[T any] struct {
	val T
	err error
}

// Map applies a function to the value of the Opt, unless the Opt is empty.
// If you need to change the type inside the Opt you will have to use Map.
func (o Opt[T]) Map(f func(T) T) Opt[T] {
	if o.Empty() {
		return o
	}
	return Of(f(o.val))
}

// Must returns the value wrapped by this opt or panics if there isn't one.
// It should be considered a last resort to use this function. It is almost always
// better to use Or, OnErr, or Return instead.
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

// Empty returns true if the option holds ErrEmpty or any other error.
func (o Opt[T]) Empty() bool {
	return o.err != nil
}

// Ok returns true if there is no error associated with this opt.
// It is guaranteed to be valid to call Opt.Must() if Opt.Ok() returns true.
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
