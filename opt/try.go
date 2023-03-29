package opt

type FuncSourceErr[T any] func() (T, error)

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
