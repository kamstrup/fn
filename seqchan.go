package fn

import "github.com/kamstrup/fn/opt"

type Chan[T any] <-chan T

// ChanOf returns a Seq that reads a channel until it is closed.
func ChanOf[T any](ch <-chan T) Seq[T] {
	return Chan[T](ch)
}

func (c Chan[T]) ForEach(f Func1[T]) Seq[T] {
	for t := range c {
		f(t)
	}

	return SeqEmpty[T]()
}

func (c Chan[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	i := 0
	for t := range c {
		f(i, t)
		i++
	}

	return SeqEmpty[T]()
}

func (c Chan[T]) Len() (int, bool) {
	return LenUnknown, false
}

func (c Chan[T]) Array() Array[T] {
	return Into(nil, Append[T], c.Seq()).Or(nil) // careful: errors silently dropped
}

func (c Chan[T]) Take(n int) (Array[T], Seq[T]) {
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
	return head, SeqEmpty[T]()
}

func (c Chan[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	var arr []T
	for t := range c {
		if pred(t) {
			arr = append(arr, t)
		} else {
			return arr, ConcatOf(SingletOf(t), c.Seq())
		}
	}

	// if we get here c was closed
	return arr, SeqEmpty[T]()
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

	return SeqEmpty[T]() // if we get here c was closed
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
		return opt.Empty[T](), SeqEmpty[T]()
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
