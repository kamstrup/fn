package seq

import "github.com/kamstrup/fn/opt"

// Set represents a collection of unique elements, represented as a standard map of empty structs.
// Sets can be used directly as go maps if you instantiate them via SetAs() or as a seq.Set literal.
// This means that you can use indexing with K, call len(set), and mutate a Set.
//
// Important: Set, as all go maps, do not have an intrinsic sort order. Methods
// returning a subset of the elements will return a random sample. Methods with
// this caveat include Seq.Take, Seq.TakeWhile, Seq.Skip, and Seq.First.
//
// # Examples:
//
//	// Sets can be created as literals
//	mySet := seq.Set[string]{"one": {}, "two": {}}
//
//	// They can be created with make()
//	emptySetWithCap10 := make(seq.Set[string], 10)
//
//	// You can call len()
//	fmt.Println("Length of mySet:", len(mySet))
//
//	// You can iterate with an idiomatic for-loop
//	for k := range mySet { fmt.Println("Key:", k) }
//
//	// You can check for element presence
//	_, hasTwo := mySet["two"]
//	hasTwoAlt := mySet.Contains("two")
type Set[K comparable] map[K]struct{}

// SetOf returns a Seq representation of standard Go set.
// Sets can be used directly as go maps if you instantiate them via SetAs().
//
// Important: Set, as all go maps, do not have an intrinsic sort order. Methods
// returning a subset of the elements will return a random sample. Methods with
// this caveat include Seq.Take, Seq.TakeWhile, Seq.Skip, and Seq.First.
func SetOf[K comparable](s map[K]struct{}) Seq[K] {
	// NOTE: Ideally this function would return Set[K]
	// and the compiler would infer that this is a valid Seq[K].
	// Alas, as of Go 1.19 this is not possible.
	// See https://github.com/golang/go/issues/41176
	return Set[K](s)
}

// SetOfArgs returns a variable argument list as a Seq.
// If you need to do set operations on the return value you can use SetAsArgs.
func SetOfArgs[K comparable](ks ...K) Seq[K] {
	return SetAsArgs(ks...)
}

// SetAs returns a Set. You can cast the set to a Seq by calling Set.Seq().
// The Go compiler can not do the type inference required to use a Set as a Seq.
func SetAs[K comparable](s map[K]struct{}) Set[K] {
	return s
}

// SetAsArgs returns a variable argument list as a Set.
// You can cast the set to a Seq by calling Set.Seq().
// The Go compiler can not do the type inference required to use a Set as a Seq.
func SetAsArgs[K comparable](ks ...K) Set[K] {
	s := Set[K]{}
	for _, k := range ks {
		s[k] = struct{}{}
	}
	return s
}

// Seq casts the Set into a Seq. This is sometimes required because
// the Go compiler can not do the type inference required to use a Set[K] as a Seq[K].
func (s Set[K]) Seq() Seq[K] {
	return s
}

func (s Set[K]) ForEach(f Func1[K]) opt.Opt[K] {
	for k := range s {
		f(k)
	}

	return opt.Zero[K]()
}

func (s Set[K]) ForEachIndex(f Func2[int, K]) opt.Opt[K] {
	idx := 0
	for k := range s {
		f(idx, k)
		idx++
	}

	return opt.Zero[K]()
}

func (s Set[K]) Len() (int, bool) {
	return len(s), true
}

func (s Set[K]) ToSlice() Slice[K] {
	sz := len(s)
	if sz == 0 {
		return Slice[K](nil)
	}

	arr := make([]K, sz)
	idx := 0
	for k := range s {
		arr[idx] = k
		idx++
	}

	return arr
}

func (s Set[T]) Limit(n int) Seq[T] {
	return LimitOf[T](s, n)
}

func (s Set[K]) Take(n int) (Slice[K], Seq[K]) {
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

	return head, SliceOf(tail)
}

func (s Set[K]) TakeWhile(predicate Predicate[K]) (Slice[K], Seq[K]) {
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

	return head, SliceOf(tail)
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
		return Empty[K]()
	} else {
		tail = make([]K, sz-n)
	}

	for k := range s {
		if idx >= n {
			tail[idx-n] = k
		}
		idx++
	}

	return SliceOf(tail)
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

// Contains return true iff the element k is in the set
func (s Set[K]) Contains(k K) bool {
	_, ok := s[k]
	return ok
}

// Union returns a lazy seq enumerating the elements in the union of 2 sets
func (s Set[K]) Union(other Set[K]) Seq[K] {
	// return seq with the smallest number of lookups
	if len(s) >= len(other) {
		return ConcatOf[K](s, other.Where(Not(s.Contains)))
	}
	return ConcatOf[K](other, s.Where(Not(other.Contains)))
}

// Intersect returns a lazy seq enumerating the elements in the intersection of 2 sets
func (s Set[K]) Intersect(other Set[K]) Seq[K] {
	// return seq with the smallest number of lookups
	if len(s) <= len(other) {
		return s.Where(other.Contains)
	}
	return other.Where(s.Contains)
}

// Copy returns a copy of this set
func (s Set[K]) Copy() Set[K] {
	dup := make(Set[K], len(s))
	for k := range s {
		dup[k] = struct{}{}
	}
	return dup
}
