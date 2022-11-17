package fn

func SeqMap[S, T comparable](seq Seq[S], fm FuncMap[S, T]) Seq[T] {
	return mappedSeq[S, T]{
		f:   fm,
		seq: seq,
	}
}

type mappedSeq[S, T comparable] struct {
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

func (m mappedSeq[S, T]) Take(n int) (Seq[T], Seq[T]) {
	head, tail := m.seq.Take(n)
	return SeqMap(head, m.f), SeqMap(tail, m.f)
}

func (m mappedSeq[S, T]) First() (Opt[T], Seq[T]) {
	s, tail := m.seq.First()
	return OptMap(s, m.f), SeqMap(tail, m.f)
}
