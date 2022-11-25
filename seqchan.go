package fn

type chanSeq[T any] <-chan T

// ChanOf returns a Seq that reads a channel until it is closed.
func ChanOf[T any](ch <-chan T) Seq[T] {
	return chanSeq[T](ch)
}

func (c chanSeq[T]) ForEach(f Func1[T]) {
	for t := range c {
		f(t)
	}
}

func (c chanSeq[T]) ForEachIndex(f Func2[int, T]) {
	i := 0
	for t := range c {
		f(i, t)
		i++
	}
}

func (c chanSeq[T]) Len() (int, bool) {
	return LenUnknown, false
}

func (c chanSeq[T]) Array() Array[T] {
	return Into(nil, Append[T], c.seq())
}

func (c chanSeq[T]) Take(n int) (Array[T], Seq[T]) {
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

func (c chanSeq[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	var arr []T
	for t := range c {
		if pred(t) {
			arr = append(arr, t)
		} else {
			return arr, ConcatOf(ArrayOfArgs(t).Seq(), c.seq())
		}
	}

	// if we get here c was closed
	return arr, SeqEmpty[T]()
}

func (c chanSeq[T]) Skip(n int) Seq[T] {
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

func (c chanSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  c,
		pred: pred,
	}
}

func (c chanSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  c,
		pred: pred,
	}
}

func (c chanSeq[T]) First() (Opt[T], Seq[T]) {
	t, ok := <-c
	if !ok {
		return OptEmpty[T](), SeqEmpty[T]()
	}
	return OptOf(t), c
}

func (c chanSeq[T]) seq() Seq[T] {
	return c
}
