package opt

// Map converts an option into some other type.
// If you want to keep the same type it may be easier to use Opt.Map.
// If the opt is empty or an error the mapping function will not be called.
func Map[S, T any](opt Opt[S], f func(S) T) Opt[T] {
	if opt.err != nil {
		return ErrorOf[T](opt.err)
	}
	return Of(f(opt.val))
}

// Ok is just a different way of calling Opt.Ok.
// It can sometimes make seq expression read a bit easier.
func Ok[T any](opt Opt[T]) bool {
	return opt.Ok()
}
