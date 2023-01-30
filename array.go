package fn

import (
	"math/rand"
	"sort"

	"github.com/kamstrup/fn/opt"
)

type Array[T any] []T

func SeqEmpty[T any]() Seq[T] {
	return Array[T](nil)
}

// ArrayOf returns a slice cast as a Seq implemented by Array.
// This method returns a Seq instead of Array because the Go compiler
// can not do the required type inference for concrete types implementing
// generic interfaces. If you need to explicitly use an Array then you
// can use ArrayAsArgs, ArrayAs, or straight type conversion Array[T](mySlice).
func ArrayOf[T any](tt []T) Seq[T] {
	// NOTE: Ideally this function would return Array[T]
	// and the compiler would infer that this is a valid Seq[T].
	// Alas, as of Go 1.19 this is not possible.
	// See https://github.com/golang/go/issues/41176
	if len(tt) == 0 {
		return SeqEmpty[T]()
	} else if len(tt) == 1 {
		return SingletOf(tt[0])
	}
	return Array[T](tt)
}

// ArrayOfArgs is a helper for creating a Seq from a variable list of arguments.
func ArrayOfArgs[T any](tt ...T) Seq[T] {
	return Array[T](tt)
}

// ArrayAs returns a slice cast into Array.
// This lets you avoid explicit type arguments that a normal type conversion would require.
func ArrayAs[T any](tt []T) Array[T] {
	return tt
}

// ArrayAsArgs return the variable argument list as an Array.
// This is sometimes needed instead of ArrayOfArgs(), when you
// need to explicitly use an Array and not any Seq.
func ArrayAsArgs[T any](tt ...T) Array[T] {
	return tt
}

func (a Array[T]) Seq() Seq[T] {
	return a
}

func (a Array[T]) ForEach(f Func1[T]) Seq[T] {
	for _, v := range a {
		f(v)
	}

	return SeqEmpty[T]()
}

func (a Array[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	for i, v := range a {
		f(i, v)
	}

	return SeqEmpty[T]()
}

func (a Array[T]) Len() (int, bool) {
	return len(a), true
}

func (a Array[T]) Array() Array[T] {
	return a
}

func (a Array[T]) Take(n int) (Array[T], Seq[T]) {
	if len(a) <= n {
		return a, SeqEmpty[T]()
	}
	return a[:n], a[n:]
}

func (a Array[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	for i, v := range a {
		if !pred(v) {
			return a[:i], a[i:]
		}
	}
	return a, SeqEmpty[T]()
}

func (a Array[T]) Skip(n int) Seq[T] {
	if n < 0 {
		panic("must skip >= 0 elements")
	}

	if len(a) <= n {
		return SeqEmpty[T]()
	}
	return a[n:]
}

func (a Array[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Array[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Array[T]) First() (opt.Opt[T], Seq[T]) {
	if len(a) == 0 {
		return opt.Empty[T](), a
	}
	return opt.Of(a[0]), a[1:]
}

func (a Array[T]) Map(shaper FuncMap[T, T]) Seq[T] {
	return mappedSeq[T, T]{
		f:   shaper,
		seq: a,
	}
}

// Sort is special for Array Seqs since it is done in place.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
//
// Typical less-functions to use are OrderAsc, OrderDesc, OrderTupleAsc, and OrderTupleDesc.
func (a Array[T]) Sort(less FuncLess[T]) Array[T] {
	sort.Slice(a, func(i, j int) bool {
		return less(a[i], a[j])
	})
	return a
}

// Reverse is special for Array Seqs since it is done in place.
// Returns the array receiver again for easy chaining.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Array[T]) Reverse() Seq[T] {
	// TODO: This could be done as a lazy seq
	end := len(a) / 2
	for i := 0; i < end; i++ {
		swapIdx := len(a) - 1 - i
		a[i], a[swapIdx] = a[swapIdx], a[i]
	}
	return a
}

// Shuffle pseudo-randomizes the elements in the Array in place.
// Returns the array receiver again for easy chaining.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Array[T]) Shuffle() Seq[T] {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

// AsSlice is a chainable method for casting the Array into a []T
func (a Array[T]) AsSlice() []T {
	return a
}

func (a Array[T]) Error() error {
	if len(a) == 0 {
		return nil
	}

	return Error(a[0])
}
