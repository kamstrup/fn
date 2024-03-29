package opt

type FuncSourceErr[T any] func() (T, error)

// Of creates a new opt wrapping a value.
func Of[T any](t T) Opt[T] {
	return Opt[T]{val: t}
}

// ErrorOf creates a new opt wrapping an error.
func ErrorOf[T any](err error) Opt[T] {
	return Opt[T]{err: err}
}

// Empty creates a new empty opt.
// An empty opt stores the special error ErrEmpty and will respond true to Opt.Empty() and false to Opt.Ok.
// An empty opt is not the same as a zero value. Empty opts should be thought of as a general nil, or missing value.
func Empty[T any]() Opt[T] {
	return Opt[T]{err: ErrEmpty}
}

// Zero creates an opt with the zero value for the type. This is *not* the same as an empty opt!
// A zero opt response true to Opt.Ok and does not have any error associated with it.
func Zero[T any]() Opt[T] {
	return Opt[T]{}
}

// Mapper wraps an error-returning function in a function returning an opt.
// Mapper is intended for use with slice.Mapping and seq.MappingOf and similar transformations.
// See also Caller.
//
// # Example
//
//		strInts := seq.SliceOfArgs("1", "two", "3")
//		ints := seq.MappingOf(strInts, opt.Mapper(strconv.Atoi)).
//			Where(opt.Ok[int]).
//			ToSlice()
//
//	 // ints is [opt.Of(1), opt.Of(3)]
func Mapper[S, T any](f func(S) (T, error)) func(S) Opt[T] {
	return func(s S) Opt[T] {
		t, err := f(s)
		return Returning(t, err)
	}
}

// Caller wraps an error-returning function in a function returning an opt.
// Caller is intended for use with seq.SourceOf and similar situation where
// you repeatedly call a source function that can error.
// See also Mapper.
//
// # Example
//
//	 i := 0
//		ints, tail := seq.SourceOf(opt.Caller(func() (int, error) {
//			i++
//			if i > 3 {
//				return 0, numTooBigError
//			}
//			return i, nil
//		})).TakeWhile(opt.Ok[int])
//
//	 // ints = [opt.Of(1), opt.Of(2), opt.Of(3)]
//	 // tail wraps numTooBigError
func Caller[T any](f FuncSourceErr[T]) func() Opt[T] {
	return func() Opt[T] {
		t, err := f()
		return Returning(t, err)
	}
}

// Returning creates an option from a value and an error.
//
// Example:
//
//	fileOpt := opt.Returning(os.Open("/tmp/README.txt"))
func Returning[T any](t T, err error) Opt[T] {
	if err != nil {
		return ErrorOf[T](err)
	}
	return Of(t)
}

// Recovering a function returning an opt with the result.
// If the function panics it is recovered and returned as ErrPanic.
func Recovering[T any](f FuncSourceErr[T]) (op Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			op = ErrorOf[T](ErrPanic{V: r})
		}
	}()

	return Returning(f())
}

// RecoveringMapper is like Mapper but recovers panics produced by f.
func RecoveringMapper[S any, T any](f func(S) (T, error)) func(S) Opt[T] {
	return func(s S) (op Opt[T]) {
		defer func() {
			if r := recover(); r != nil {
				op = ErrorOf[T](ErrPanic{V: r})
			}
		}()

		return Returning(f(s))
	}
}
