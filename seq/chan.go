package seq

import "github.com/kamstrup/fn/opt"

// Chan is type wrapper for standard Go channels.
// This means that standard Go builtins and syntax applies to seq.Chan as well.
// You can fx. call len() and cap() on them, make() them, and use select- and for-range expressions.
type Chan[T any] chan T

// ChanOf returns a Seq that reads a channel until it is closed.
func ChanOf[T any](ch chan T) Seq[T] {
	return Chan[T](ch)
}

func (c Chan[T]) ForEach(f Func1[T]) opt.Opt[T] {
	for t := range c {
		f(t)
	}

	return opt.Zero[T]()
}

func (c Chan[T]) ForEachIndex(f Func2[int, T]) opt.Opt[T] {
	i := 0
	for t := range c {
		f(i, t)
		i++
	}

	return opt.Zero[T]()
}

func (c Chan[T]) Len() (int, bool) {
	return LenUnknown, false
}

func (c Chan[T]) ToSlice() Slice[T] {
	return Reduce(MakeSlice[T], nil, c.Seq()).Or(nil) // careful: errors silently dropped
}

func (c Chan[T]) Limit(n int) Seq[T] {
	return LimitOf[T](c, n)
}

func (c Chan[T]) Take(n int) (Slice[T], Seq[T]) {
	if n == 0 {
		return []T{}, c
	}

	head := make([]T, 0, n)
	i := 0
	for t := range c {
		head = append(head, t)
		i++
		if i >= n {
			return head, c
		}
	}

	// If we get here c was closed
	return head, Empty[T]()
}

func (c Chan[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	var arr []T
	for t := range c {
		if pred(t) {
			arr = append(arr, t)
		} else {
			return arr, PrependOf(t, c.Seq())
		}
	}

	// if we get here c was closed
	return arr, Empty[T]()
}

func (c Chan[T]) Skip(n int) Seq[T] {
	if n == 0 {
		return c
	}

	i := 0
	for _ = range c {
		i++
		if i >= n {
			return c
		}
	}

	return Empty[T]() // if we get here c was closed
}

func (c Chan[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  c,
		pred: pred,
	}
}

func (c Chan[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  c,
		pred: pred,
	}
}

func (c Chan[T]) First() (opt.Opt[T], Seq[T]) {
	t, ok := <-c
	if !ok {
		return opt.Empty[T](), Empty[T]()
	}
	return opt.Of(t), c
}

func (c Chan[T]) Map(shaper FuncMap[T, T]) Seq[T] {
	return mappedSeq[T, T]{
		f:   shaper,
		seq: c,
	}
}

func (c Chan[T]) Seq() Seq[T] {
	return c
}
