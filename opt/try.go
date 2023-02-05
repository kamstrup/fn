package opt

type FuncSourceErr[T any] func() (T, error)

// Try wraps an error-returning function in a function returning an opt.
// Try is intended for use with slice.Mapping and seq.MappingOf and similar transformations.
//
// # Example
//
//	filenames := seq.SliceOfArgs("file1.txt", "file2.txt")
//	openFiles := seq.MappingOf(filenames, opt.Try(os.Open)).
//	       Where(opt.Ok)
//	fileContents := se.MappingOf(openFiles, io.ReadAll)
func Try[S, T any](f func(S) (T, error)) func(S) Opt[T] {
	return func(s S) Opt[T] {
		t, err := f(s)
		return Returning(t, err)
	}
}

// Returning creates an option from a value and an error.
func Returning[T any](t T, err error) Opt[T] {
	if err != nil {
		return ErrorOf[T](err)
	}
	return Of(t)
}

// Call a function returning an option with the result
func Call[T any](f FuncSourceErr[T]) Opt[T] {
	t, err := f()
	return Returning(t, err)
}

// CallRecover a function returning an opt with the result.
// If the function panics it is recovered and returned as ErrPanic.
func CallRecover[T any](f FuncSourceErr[T]) (op Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			op = ErrorOf[T](ErrPanic{V: r})
		}
	}()

	return Call(f)
}

// Apply calls a function with an argument and returns the result wrapped in an opt.
func Apply[S any, T any](f func(S) (T, error), s S) Opt[T] {
	t, err := f(s)
	return Returning(t, err)
}

// ApplyRecover calls a function with an argument and returns the result wrapped in an opt.
// If the function panics it is recovered and returned as ErrPanic.
func ApplyRecover[S any, T any](f func(S) (T, error), s S) (op Opt[T]) {
	defer func() {
		if r := recover(); r != nil {
			op = ErrorOf[T](ErrPanic{V: r})
		}
	}()

	return Apply(f, s)
}
