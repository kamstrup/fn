package fn

import "github.com/kamstrup/fn/opt"

// Set represents a collection of unique elements, represented as a standard map of empty structs.
type Set[K comparable] map[K]struct{}

// SetOf returns a Seq representation of standard Go set.
// Sets can be used directly as Go maps if you instantiate them via SetAs().
func SetOf[K comparable](s map[K]struct{}) Seq[K] {
	// NOTE: Ideally this function would return Set[K]
	// and the compiler would infer that this is a valid Seq[K].
	// Alas, as of Go 1.19 this is not possible.
	// See https://github.com/golang/go/issues/41176
	return Set[K](s)
}

// SetAs returns a Set. You can cast the set to a Seq by calling Set.Seq().
// The Go compiler can not do the type inference required to use a Set as a Seq.
func SetAs[K comparable](s map[K]struct{}) Set[K] {
	return s
}

// Seq casts the Set into a Seq. This is sometimes required because
// the Go compiler can not do the type inference required to use a Set[K] as a Seq[K].
func (s Set[K]) Seq() Seq[K] {
	return s
}

func (s Set[K]) ForEach(f Func1[K]) Seq[K] {
	for k := range s {
		f(k)
	}

	return SeqEmpty[K]()
}

func (s Set[K]) ForEachIndex(f Func2[int, K]) Seq[K] {
	idx := 0
	for k := range s {
		f(idx, k)
		idx++
	}

	return SeqEmpty[K]()
}

func (s Set[K]) Len() (int, bool) {
	return len(s), true
}

func (s Set[K]) Array() Array[K] {
	sz := len(s)
	if sz == 0 {
		return Array[K](nil)
	}

	arr := make([]K, sz)
	idx := 0
	for k := range s {
		arr[idx] = k
		idx++
	}

	return arr
}

func (s Set[K]) Take(n int) (Array[K], Seq[K]) {
	// Taking the "first n elements" from a map[K]V does *almost* never make sense,
	// since maps in Go a deliberately not ordered consistently.
	// We provide the feature for completeness.

	if n == 0 {
		return []K{}, s
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

	return head, ArrayOf(tail)
}

func (s Set[K]) TakeWhile(predicate Predicate[K]) (Array[K], Seq[K]) {
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

	return head, ArrayOf(tail)
}

func (s Set[K]) Skip(n int) Seq[K] {
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

func (s Set[K]) Where(p Predicate[K]) Seq[K] {
	return whereSeq[K]{
		seq:  s,
		pred: p,
	}
}

func (s Set[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  s,
		pred: pred,
	}
}

func (s Set[K]) First() (opt.Opt[K], Seq[K]) {
	head, tail := s.Take(1)
	first, _ := head.First()
	return first, tail
}

func (s Set[K]) Map(shaper FuncMap[K, K]) Seq[K] {
	return mappedSeq[K, K]{
		f:   shaper,
		seq: s,
	}
}
