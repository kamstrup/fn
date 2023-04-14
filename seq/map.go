package seq

import "github.com/kamstrup/fn/opt"

// Map is a type wrapper for Go maps exposing them as a Seq of Tuple[K,V].
//
// A Map can be used directly as a go map if you instantiate them via MapAs() or as a seq.Map literal.
// That means you can index by K and call len() and cap() on it.
//
// Important: Map, as all go maps, do not have an intrinsic sort order. Methods
// returning a subset of the elements will return a random sample. Methods with
// this caveat include Seq.Take, Seq.TakeWhile, Seq.Skip, and Seq.First.
//
// # Examples:
//
//	// Maps can be created as literals
//	myMap := seq.Map[string]int{"one": 1, "two: 2}
//
//	// They can be created with make()
//	emptyMapWithCap10 := make(seq.Map[string,int], 10)
//
//	// You can call len()
//	fmt.Println("Length of myMap:", len(myMap))
//
//	// You can iterate with an idiomatic for-loop
//	for k, v := range myMap { fmt.Println("Key:", k, "Value:", v) }
//
//	// You can access elements
//	twoInt := myMap["two"]
type Map[K comparable, V any] map[K]V

// MapOf returns a map cast as a Seq implemented by Map.
// If you need to explicitly use a Map then please use MapAs
// and call Map.Seq() when you need to use it as a Seq.
// The go compiler can not do the type inference required to use
// a Map as a Seq.
//
// If you are looking for a way to convert a seq via a mapping operation
// please look at MappingOf.
//
// Important: Map, as all go maps, do not have an intrinsic sort order. Methods
// returning a subset of the elements will return a random sample. Methods with
// this caveat include Seq.Take, Seq.TakeWhile, Seq.Skip, and Seq.First.
func MapOf[K comparable, V any](m map[K]V) Seq[Tuple[K, V]] {
	// NOTE: Ideally this function would return Map[K,V]
	// and the compiler would infer that this is a valid Seq[Tuple[K, V]].
	// Alas, as of Go 1.19 this is not possible.
	// See https://github.com/golang/go/issues/41176
	return Map[K, V](m)
}

// MapAs returns a map cast as an Map. To use an Map as a Seq you can call Map.Seq().
// This is sometimes needed because the Go compiler can not do the type
// inference required to use an Map as a Seq.
// Since an Map is a Go map you can use normal indexing to access elements.
func MapAs[K comparable, V any](m map[K]V) Map[K, V] {
	return m
}

// Keys returns a seq over the keys in the map.
func (a Map[K, V]) Keys() Seq[K] {
	return MappingOf(a.Seq(), TupleKey[K, V])
}

// Values returns a seq over the values in the map.
func (a Map[K, V]) Values() Seq[V] {
	return MappingOf(a.Seq(), TupleValue[K, V])
}

// Seq returns the Map cast as a Seq.
func (a Map[K, V]) Seq() Seq[Tuple[K, V]] {
	return a
}

func (a Map[K, V]) ForEach(f Func1[Tuple[K, V]]) Seq[Tuple[K, V]] {
	for k, v := range a {
		f(Tuple[K, V]{k, v})
	}

	return Empty[Tuple[K, V]]()
}

func (a Map[K, V]) ForEachIndex(f Func2[int, Tuple[K, V]]) Seq[Tuple[K, V]] {
	idx := 0
	for k, v := range a {
		f(idx, Tuple[K, V]{k, v})
		idx++
	}

	return Empty[Tuple[K, V]]()
}

func (a Map[K, V]) Len() (int, bool) {
	return len(a), true
}

func (a Map[K, V]) ToSlice() Slice[Tuple[K, V]] {
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

func (a Map[K, V]) Limit(n int) Seq[Tuple[K, V]] {
	return LimitOf[Tuple[K, V]](a, n)
}

func (a Map[K, V]) Take(n int) (Slice[Tuple[K, V]], Seq[Tuple[K, V]]) {
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

func (a Map[K, V]) TakeWhile(predicate Predicate[Tuple[K, V]]) (Slice[Tuple[K, V]], Seq[Tuple[K, V]]) {
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

func (a Map[K, V]) Skip(n int) Seq[Tuple[K, V]] {
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
		return Empty[Tuple[K, V]]()
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

func (a Map[K, V]) Where(p Predicate[Tuple[K, V]]) Seq[Tuple[K, V]] {
	return whereSeq[Tuple[K, V]]{
		seq:  a,
		pred: p,
	}
}

func (a Map[K, V]) While(pred Predicate[Tuple[K, V]]) Seq[Tuple[K, V]] {
	return whileSeq[Tuple[K, V]]{
		seq:  a,
		pred: pred,
	}
}

func (a Map[K, V]) First() (opt.Opt[Tuple[K, V]], Seq[Tuple[K, V]]) {
	head, tail := a.Take(1)
	first, _ := head.First()
	return first, tail
}

func (a Map[K, V]) Map(shaper FuncMap[Tuple[K, V], Tuple[K, V]]) Seq[Tuple[K, V]] {
	return mappedSeq[Tuple[K, V], Tuple[K, V]]{
		f:   shaper,
		seq: a,
	}
}

func (a Map[K, V]) Contains(k K) bool {
	_, ok := a[k]
	return ok
}

func (a Map[K, V]) Get(k K) opt.Opt[V] {
	if v, ok := a[k]; ok {
		return opt.Of(v)
	}
	return opt.Empty[V]()
}
