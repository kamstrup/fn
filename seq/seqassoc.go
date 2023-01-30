package seq

import "github.com/kamstrup/fn/opt"

// Assoc is a type wrapper for Go maps exposing them as a Seq of Tuple[K,V].
//
// An Assoc can be used directly as a go map if you instantiate them via AssocAs().
//
// Important: Assoc, as all go maps, do not have an intrinsic sort order. Methods
// returning a subset of the elements will return a random sample. Methods with
// this caveat include Seq.Take, Seq.TakeWhile, Seq.Skip, and Seq.First.
type Assoc[K comparable, V any] map[K]V

// AssocOf returns a map cast as a Seq implemented by Assoc.
// If you need to explicitly use an Assoc then please use AssocAs
// and call Assoc.Seq() when you need to use it as a Seq.
// The Go compiler can not do the type inference required to use
// an Assoc as a Seq.
//
// Important: Assoc, as all go maps, do not have an intrinsic sort order. Methods
// returning a subset of the elements will return a random sample. Methods with
// this caveat include Seq.Take, Seq.TakeWhile, Seq.Skip, and Seq.First.
func AssocOf[K comparable, V any](m map[K]V) Seq[Tuple[K, V]] {
	// NOTE: Ideally this function would return Assoc[K,V]
	// and the compiler would infer that this is a valid Seq[Tuple[K, V]].
	// Alas, as of Go 1.19 this is not possible.
	// See https://github.com/golang/go/issues/41176
	return Assoc[K, V](m)
}

// AssocAs returns a map cast as an Assoc. To use an Assoc as a Seq you can call Assoc.Seq().
// This is sometimes needed because the Go compiler can not do the type
// inference required to use an Assoc as a Seq.
// Since an Assoc is a Go map you can use normal indexing to access elements.
func AssocAs[K comparable, V any](m map[K]V) Assoc[K, V] {
	return m
}

// Seq returns the Assoc cast as a Seq.
func (a Assoc[K, V]) Seq() Seq[Tuple[K, V]] {
	return a
}

func (a Assoc[K, V]) ForEach(f Func1[Tuple[K, V]]) Seq[Tuple[K, V]] {
	for k, v := range a {
		f(Tuple[K, V]{k, v})
	}

	return SeqEmpty[Tuple[K, V]]()
}

func (a Assoc[K, V]) ForEachIndex(f Func2[int, Tuple[K, V]]) Seq[Tuple[K, V]] {
	idx := 0
	for k, v := range a {
		f(idx, Tuple[K, V]{k, v})
		idx++
	}

	return SeqEmpty[Tuple[K, V]]()
}

func (a Assoc[K, V]) Len() (int, bool) {
	return len(a), true
}

func (a Assoc[K, V]) Values() Slice[Tuple[K, V]] {
	sz := len(a)
	if sz == 0 {
		return Slice[Tuple[K, V]](nil)
	}

	arr := make([]Tuple[K, V], sz)
	idx := 0
	for k, v := range a {
		arr[idx] = Tuple[K, V]{k, v}
		idx++
	}

	return arr
}

func (a Assoc[K, V]) Take(n int) (Slice[Tuple[K, V]], Seq[Tuple[K, V]]) {
	// Taking the "first n elements" from a map[K]V does *almost* never make sense,
	// since maps in Go a deliberately not ordered consistently.
	// We provide the feature for completeness.

	if n == 0 {
		return []Tuple[K, V]{}, a
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

	return head, SliceOf(tail)
}

func (a Assoc[K, V]) TakeWhile(predicate Predicate[Tuple[K, V]]) (Slice[Tuple[K, V]], Seq[Tuple[K, V]]) {
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

	return head, SliceOf(tail)
}

func (a Assoc[K, V]) Skip(n int) Seq[Tuple[K, V]] {
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

	return SliceOf(tail)
}

func (a Assoc[K, V]) Where(p Predicate[Tuple[K, V]]) Seq[Tuple[K, V]] {
	return whereSeq[Tuple[K, V]]{
		seq:  a,
		pred: p,
	}
}

func (a Assoc[K, V]) While(pred Predicate[Tuple[K, V]]) Seq[Tuple[K, V]] {
	return whileSeq[Tuple[K, V]]{
		seq:  a,
		pred: pred,
	}
}

func (a Assoc[K, V]) First() (opt.Opt[Tuple[K, V]], Seq[Tuple[K, V]]) {
	head, tail := a.Take(1)
	first, _ := head.First()
	return first, tail
}

func (a Assoc[K, V]) Map(shaper FuncMap[Tuple[K, V], Tuple[K, V]]) Seq[Tuple[K, V]] {
	return mappedSeq[Tuple[K, V], Tuple[K, V]]{
		f:   shaper,
		seq: a,
	}
}
