package fn

type Array[T any] struct {
	vals []T
}

func SeqEmpty[T any]() Seq[T] {
	return Array[T]{}
}

func ArrayOf[T any](tt []T) Array[T] {
	return Array[T]{tt}
}

func ArrayOfArgs[T any](tt ...T) Array[T] {
	return Array[T]{tt}
}

// Seq is a helper function for letting the Go compiler understand that Array[T] implements Seq[T]
func (a Array[T]) Seq() Seq[T] {
	return a
}

func (a Array[T]) ForEach(f Func1[T]) {
	for _, v := range a.vals {
		f(v)
	}
}

func (a Array[T]) ForEachIndex(f Func2[int, T]) {
	for i, v := range a.vals {
		f(i, v)
	}
}

func (a Array[T]) Len() int {
	return len(a.vals)
}

func (a Array[T]) Array() Array[T] {
	return a
}

func (a Array[T]) Take(n int) (Array[T], Seq[T]) {
	if a.Len() <= n {
		return a, SeqEmpty[T]()
	}
	return Array[T]{vals: a.vals[:n]}, Array[T]{vals: a.vals[n:]}
}

func (a Array[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	for i, v := range a.vals {
		if !pred(v) {
			return ArrayOf(a.vals[:i]), ArrayOf(a.vals[i:])
		}
	}
	return a, SeqEmpty[T]()
}

func (a Array[T]) Skip(n int) Seq[T] {
	if a.Len() <= n {
		return SeqEmpty[T]()
	}
	return Array[T]{vals: a.vals[n:]}
}

func (a Array[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Array[T]) First() (Opt[T], Seq[T]) {
	if len(a.vals) == 0 {
		return OptEmpty[T](), a
	}
	return OptOf(a.vals[0]), ArrayOf(a.vals[1:])
}
