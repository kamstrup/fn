package fn

type setSeq[K comparable] map[K]struct{}

func SetOf[K comparable](s map[K]struct{}) Seq[K] {
	return setSeq[K](s)
}

func (s setSeq[K]) ForEach(f Func1[K]) {
	for k := range s {
		f(k)
	}
}

func (s setSeq[K]) ForEachIndex(f Func2[int, K]) {
	idx := 0
	for k := range s {
		f(idx, k)
		idx++
	}
}

func (s setSeq[K]) Len() int {
	return len(s)
}

func (s setSeq[K]) Array() Array[K] {
	sz := len(s)
	if sz == 0 {
		return ArrayOf[K](nil)
	}

	arr := make([]K, sz)
	idx := 0
	for k := range s {
		arr[idx] = k
		idx++
	}

	return ArrayOf(arr)
}

func (s setSeq[K]) Take(n int) (Array[K], Seq[K]) {
	// Taking the "first n elements" from a map[K]V does *almost* never make sense,
	// since maps in Go a deliberately not ordered consistently.
	// We provide the feature for completeness.

	if n == 0 {
		return ArrayOf([]K{}), s
	}

	var (
		head []K
		tail []K
		idx  int
	)
	sz := len(s)
	if n >= sz {
		head = make([]K, sz)
		// tail will be empty, we do not have n elements
	} else {
		head = make([]K, n)
		tail = make([]K, sz-n)
	}

	for k := range s {
		if idx >= n {
			tail[idx-n] = k
		} else {
			head[idx] = k
		}
		idx++
	}

	return ArrayOf(head), ArrayOf(tail)
}

func (s setSeq[K]) TakeWhile(predicate Predicate[K]) (Array[K], Seq[K]) {
	// TakeWhile makes a *little* more sense on a map[K]V than Take(n) does,
	// but not much... For the rare case where someone needs it we provide the feature for completeness.
	// Example: Collect up to N random values from the map where V has some property.
	var (
		// TODO: We could be memory efficient here and have head+tail share a slice of size len(a.m)
		// ... but maybe bad for GC, since head or tail can not be GCed individually anymore
		head []K
		tail []K
	)

	for k := range s {
		if len(tail) > 0 { // after first time predicate(t) is false, don't call it again
			tail = append(tail, k)
		} else if predicate(k) {
			head = append(head, k)
		} else {
			tail = append(tail, k)
		}
	}

	return ArrayOf(head), ArrayOf(tail)
}

func (s setSeq[K]) Skip(n int) Seq[K] {
	// Skipping the "first n elements" from a map[K]V does *almost* never make sense,
	// since maps in Go a deliberately not ordered consistently.
	// We provide the feature for completeness.
	if n == 0 {
		return s
	}

	var (
		tail []K
		idx  int
	)
	sz := len(s)
	if n >= sz {
		return SeqEmpty[K]()
	} else {
		tail = make([]K, sz-n)
	}

	for k := range s {
		if idx >= n {
			tail[idx-n] = k
		}
		idx++
	}

	return ArrayOf(tail)
}

func (s setSeq[K]) Where(p Predicate[K]) Seq[K] {
	return whereSeq[K]{
		seq:  s,
		pred: p,
	}
}

func (s setSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  s,
		pred: pred,
	}
}

func (s setSeq[K]) First() (Opt[K], Seq[K]) {
	head, tail := s.Take(1)
	first, _ := head.First()
	return first, tail
}
