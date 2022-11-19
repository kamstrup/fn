package fn

var _ Seq[Tuple[int, int]] = assocSeq[int, int]{}

type assocSeq[K comparable, V any] map[K]V

func AssocOf[K comparable, V any](m map[K]V) Seq[Tuple[K, V]] {
	return assocSeq[K, V](m)
}

func (a assocSeq[K, V]) ForEach(f Func1[Tuple[K, V]]) {
	for k, v := range a {
		f(Tuple[K, V]{k, v})
	}
}

func (a assocSeq[K, V]) ForEachIndex(f Func2[int, Tuple[K, V]]) {
	idx := 0
	for k, v := range a {
		f(idx, Tuple[K, V]{k, v})
		idx++
	}
}

func (a assocSeq[K, V]) Len() int {
	return len(a)
}

func (a assocSeq[K, V]) Array() Array[Tuple[K, V]] {
	sz := len(a)
	if sz == 0 {
		return ArrayOf[Tuple[K, V]](nil)
	}

	arr := make([]Tuple[K, V], sz)
	idx := 0
	for k, v := range a {
		arr[idx] = Tuple[K, V]{k, v}
		idx++
	}

	return ArrayOf(arr)
}

func (a assocSeq[K, V]) Take(n int) (Array[Tuple[K, V]], Seq[Tuple[K, V]]) {
	// Taking the "first n elements" from a map[K]V does *almost* never make sense,
	// since maps in Go a deliberately not ordered consistently.
	// We provide the feature for completeness.

	if n == 0 {
		return ArrayOf([]Tuple[K, V]{}), a
	}

	var (
		head []Tuple[K, V]
		tail []Tuple[K, V]
		idx  int
	)
	sz := len(a)
	if n >= sz {
		head = make([]Tuple[K, V], sz)
		// tail will be empty, we do not have n elements
	} else {
		head = make([]Tuple[K, V], n)
		tail = make([]Tuple[K, V], sz-n)
	}

	for k, v := range a {
		if idx >= n {
			tail[idx-n] = Tuple[K, V]{k, v}
		} else {
			head[idx] = Tuple[K, V]{k, v}
		}
		idx++
	}

	return ArrayOf(head), ArrayOf(tail)
}

func (a assocSeq[K, V]) TakeWhile(predicate Predicate[Tuple[K, V]]) (Array[Tuple[K, V]], Seq[Tuple[K, V]]) {
	// TakeWhile makes a *little* more sense on a map[K]V than Take(n) does,
	// but not much... For the rare case where someone needs it we provide the feature for completeness.
	// Example: Collect up to N random values from the map where V has some property.
	var (
		// TODO: We could be memory efficient here and have head+tail share a slice of size len(a.m)
		// ... but maybe bad for GC, since head or tail can not be GCed individually anymore
		head []Tuple[K, V]
		tail []Tuple[K, V]
	)

	for k, v := range a {
		t := Tuple[K, V]{k, v}
		if len(tail) > 0 { // after first time predicate(t) is false, don't call it again
			tail = append(tail, t)
		} else if predicate(t) {
			head = append(head, t)
		} else {
			tail = append(tail, t)
		}
	}

	return ArrayOf(head), ArrayOf(tail)
}

func (a assocSeq[K, V]) Skip(n int) Seq[Tuple[K, V]] {
	// Skipping the "first n elements" from a map[K]V does *almost* never make sense,
	// since maps in Go a deliberately not ordered consistently.
	// We provide the feature for completeness.
	if n == 0 {
		return a
	}

	var (
		tail []Tuple[K, V]
		idx  int
	)
	sz := len(a)
	if n >= sz {
		return SeqEmpty[Tuple[K, V]]()
	} else {
		tail = make([]Tuple[K, V], sz-n)
	}

	for k, v := range a {
		if idx >= n {
			tail[idx-n] = Tuple[K, V]{k, v}
		}
		idx++
	}

	return ArrayOf(tail)
}

func (a assocSeq[K, V]) Where(p Predicate[Tuple[K, V]]) Seq[Tuple[K, V]] {
	return whereSeq[Tuple[K, V]]{
		seq:  a,
		pred: p,
	}
}

func (a assocSeq[K, V]) While(pred Predicate[Tuple[K, V]]) Seq[Tuple[K, V]] {
	return whileSeq[Tuple[K, V]]{
		seq:  a,
		pred: pred,
	}
}

func (a assocSeq[K, V]) First() (Opt[Tuple[K, V]], Seq[Tuple[K, V]]) {
	head, tail := a.Take(1)
	first, _ := head.First()
	return first, tail
}
