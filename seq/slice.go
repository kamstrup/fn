package seq

import (
	"math/rand"
	"sort"

	"github.com/kamstrup/fn/opt"
)

// Slice is a type wrapper for standard go slices.
// You can either create your slices via normal type conversion, Slice[T](mySlice),
// or more easily (with type inference) via the static constructor SliceAs, and SliceAsArgs.
//
// You can use numeric indexing and call len(slice) directly on Slice instances,
// and any function expecting a normal slice, []T, can take a Slice[T] as well.
//
// # Examples:
//
//	// Slices can be created as literals
//	mySlice := seq.Slice[string]{"one", "two"}
//
//	// They can be allocated with make()
//	emptySliceWithCap10 := make(seq.Slice[string], 0, 10)
//
//	// You can call len()
//	fmt.Println("Length of mySlice:", len(mySlice))
//
//	// You can iterate with an idiomatic for-loop
//	for i, v := range mySlice { fmt.Println("Index:", i, "Value:", v) }
//
//	// You can access elements by index
//	twoString := mySlice[1]
type Slice[T any] []T

// Empty returns an empty seq
func Empty[T any]() Seq[T] {
	return Slice[T](nil)
}

// SliceOf returns a slice cast as a Seq implemented by Slice.
// This method returns a Seq instead of Slice because the Go compiler
// can not do the required type inference for concrete types implementing
// generic interfaces. If you need to explicitly use a Slice then you
// can use SliceAsArgs, SliceAs, or straight type conversion Slice[T](mySlice).
func SliceOf[T any](tt []T) Seq[T] {
	// NOTE: Ideally this function would return Slice[T]
	// and the compiler would infer that this is a valid Seq[T].
	// Alas, as of Go 1.19 this is not possible.
	// See https://github.com/golang/go/issues/41176
	if len(tt) == 0 {
		return Empty[T]()
	} else if len(tt) == 1 {
		return SingletOf(tt[0])
	}
	return Slice[T](tt)
}

// SliceOfArgs is a helper for creating a Seq from a variable list of arguments.
func SliceOfArgs[T any](tt ...T) Seq[T] {
	return Slice[T](tt)
}

// SliceAs returns a slice cast into Slice.
// This lets you avoid explicit type arguments that a normal type conversion would require.
func SliceAs[T any](tt []T) Slice[T] {
	return tt
}

// SliceAsArgs return the variable argument list as a Slice.
// This is sometimes needed instead of SliceOfArgs(), when you
// need to explicitly use a Slice and not any Seq.
func SliceAsArgs[T any](tt ...T) Slice[T] {
	return tt
}

func (a Slice[T]) Seq() Seq[T] {
	return a
}

func (a Slice[T]) ForEach(f Func1[T]) Seq[T] {
	for _, v := range a {
		f(v)
	}

	return Empty[T]()
}

func (a Slice[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	for i, v := range a {
		f(i, v)
	}

	return Empty[T]()
}

func (a Slice[T]) Len() (int, bool) {
	return len(a), true
}

func (a Slice[T]) ToSlice() Slice[T] {
	return a.Copy() // must return a new instance
}

func (a Slice[T]) Limit(n int) Seq[T] {
	if len(a) <= n {
		return a
	}
	return a[:n]
}

func (a Slice[T]) Take(n int) (Slice[T], Seq[T]) {
	if len(a) <= n {
		return a, Empty[T]()
	}
	return a[:n], a[n:]
}

func (a Slice[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	for i, v := range a {
		if !pred(v) {
			return a[:i], a[i:]
		}
	}
	return a, Empty[T]()
}

func (a Slice[T]) Skip(n int) Seq[T] {
	if n < 0 {
		panic("must skip >= 0 elements")
	}

	if len(a) <= n {
		return Empty[T]()
	}
	return a[n:]
}

func (a Slice[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Slice[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Slice[T]) First() (opt.Opt[T], Seq[T]) {
	if len(a) == 0 {
		return opt.Empty[T](), a
	}
	return opt.Of(a[0]), a[1:]
}

func (a Slice[T]) Map(shaper FuncMap[T, T]) Seq[T] {
	return mappedSeq[T, T]{
		f:   shaper,
		seq: a,
	}
}

// Last behaves like the Last function in the seq package, but allows for easier chaining.
func (a Slice[T]) Last() opt.Opt[T] {
	if len(a) > 0 {
		return opt.Of(a[len(a)-1])
	}
	return opt.Empty[T]()
}

// One behaves like the One function in the seq package, but allows for easier chaining.
func (a Slice[T]) One() opt.Opt[T] {
	switch len(a) {
	case 0:
		return opt.ErrorOf[T](opt.ErrEmpty)
	case 1:
		return opt.Of(a[0])
	default:
		return opt.ErrorOf[T](ErrNotOne)
	}
}

// Sort is special for Slice Seqs since it is done in place.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
//
// Typical less-functions to use are OrderAsc, OrderDesc, OrderTupleAsc, and OrderTupleDesc.
func (a Slice[T]) Sort(less FuncLess[T]) Slice[T] {
	sort.Slice(a, func(i, j int) bool {
		return less(a[i], a[j])
	})
	return a
}

// Reverse is special for Slice Seqs since it is done in place.
// Returns the array receiver again for easy chaining.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Slice[T]) Reverse() Seq[T] {
	// TODO: This could be done as a lazy seq
	end := len(a) / 2
	for i := 0; i < end; i++ {
		swapIdx := len(a) - 1 - i
		a[i], a[swapIdx] = a[swapIdx], a[i]
	}
	return a
}

// Shuffle pseudo-randomizes the elements in the Slice in place.
// Returns the array receiver again for easy chaining.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Slice[T]) Shuffle() Seq[T] {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

// Copy returns a copy of this slice
func (a Slice[T]) Copy() Slice[T] {
	if len(a) == 0 {
		return []T{}
	}

	cpy := make([]T, len(a))
	copy(cpy, a)
	return cpy
}

func (a Slice[T]) Error() error {
	if len(a) == 0 {
		return nil
	}

	return Error(a[0])
}
