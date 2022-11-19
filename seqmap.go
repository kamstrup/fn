package fn

// MapOf creates a new Seq that lazily converts the values, via a FuncMap,
// into another Seq. The returned Seq has the same Seq.Len() as the input Seq.
// If you are looking for ways to create a Seq from a Go map[K]V please
// look at AssocOf() or SetOf().
func MapOf[S, T any](seq Seq[S], fm FuncMap[S, T]) Seq[T] {
	return mappedSeq[S, T]{
		f:   fm,
		seq: seq,
	}
}

type mappedSeq[S, T any] struct {
	f   FuncMap[S, T]
	seq Seq[S]
}

func (m mappedSeq[S, T]) ForEach(f Func1[T]) {
	m.seq.ForEach(func(s S) {
		t := m.f(s)
		f(t)
	})
}

func (m mappedSeq[S, T]) ForEachIndex(f Func2[int, T]) {
	m.seq.ForEachIndex(func(i int, s S) {
		t := m.f(s)
		f(i, t)
	})
}

func (m mappedSeq[S, T]) Len() int {
	return m.seq.Len()
}

func (m mappedSeq[S, T]) Array() Array[T] {
	if sz := m.seq.Len(); sz != LenUnknown {
		arr := make([]T, sz)
		m.ForEachIndex(func(i int, t T) {
			arr[i] = t
		})
		return ArrayOf(arr)
	} else {
		var arr []T
		m.ForEach(func(t T) {
			arr = append(arr, t)
		})
		return ArrayOf(arr)
	}
}

func (m mappedSeq[S, T]) Take(n int) (Array[T], Seq[T]) {
	var (
		head Array[S]
		tail Seq[S]
	)
	head, tail = m.seq.Take(n)
	return MapOf[S, T](head, m.f).Array(), MapOf(tail, m.f)
}

func (m mappedSeq[S, T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	// TODO: does not really to alloc a slice, if we had a "pulling seq"
	var arr []T
	_, tail := m.seq.TakeWhile(func(s S) bool {
		t := m.f(s)
		if pred(t) {
			arr = append(arr, t)
			return true
		}
		return false
	})
	return ArrayOf(arr), MapOf(tail, m.f)
}

func (m mappedSeq[S, T]) Skip(n int) Seq[T] {
	tail := m.seq.Skip(n)
	return MapOf(tail, m.f)
}
func (m mappedSeq[S, T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  m,
		pred: pred,
	}
}

func (m mappedSeq[S, T]) First() (Opt[T], Seq[T]) {
	s, tail := m.seq.First()
	return OptMap(s, m.f), MapOf(tail, m.f)
}
