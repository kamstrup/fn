package seq

import (
	"bytes"
	"errors"
	"strings"

	"github.com/kamstrup/fn/constraints"
	"github.com/kamstrup/fn/opt"
)

type Func0 func()

type Func0Err func() error

type Func1[T any] func(T)

type Func1Err[T any] func(T) error

type Func2[S, T any] func(S, T)

type Func2Err[S, T any] func(S, T) error

// FuncMap is a function mapping type S to type T.
// Used with fx. MappingOf(), Map(), and TupleWithKey().
type FuncMap[S, T any] func(S) T

type FuncMapErr[S, T any] func(S) (T, error)

// FuncCollect is used to aggregate data and return the updated aggregation.
// Think of it as a generic form of the standard append() function in go.
// The order of the method arguments also follow the convention from append() and copy(),
// having the target destination as first argument.
type FuncCollect[T, E any] func(T, E) T

type FuncSource[T any] func() T

type FuncLess[T any] func(T, T) bool

type Predicate[T any] func(T) bool

type PredicateErr[T any] func(T) (bool, error)

// FuncUpdate is used to update an existing 'old' value compared to a 'new_' value, and returning the updated result.
type FuncUpdate[T any] func(old, new_ T) T

// Count is a FuncCollect for use with Reduce, that counts the number of elements it sees.
func Count[T any](into int, _ T) int {
	return into + 1
}

// Append is a FuncCollect for use with Reduce, that uses the standard Go append() function.
// This function works with nil or a pre-built slice as initial value.
func Append[T any](into []T, t T) []T {
	return append(into, t)
}

// MakeString is a FuncCollect for use with Reduce that writes strings into a bytes.Buffer.
// This function works with nil or a pre-built bytes.Buffer as initial value.
func MakeString(into *strings.Builder, s string) *strings.Builder {
	if into == nil {
		into = &strings.Builder{}
	}
	_, _ = into.WriteString(s)
	return into
}

// MakeBytes is a FuncCollect for use with Reduce that writes bytes into a bytes.Buffer.
// This function works with nil or a pre-built bytes.Buffer as initial value.
func MakeBytes(into *bytes.Buffer, b []byte) *bytes.Buffer {
	if into == nil {
		into = &bytes.Buffer{}
	}
	_, _ = into.Write(b)
	return into
}

// MakeMap is a FuncCollect that can take a Seq of Tuple elements and store them in a standard go map.
// This function works with nil or a pre-built map[K]V as initial value.
func MakeMap[K comparable, V any](into map[K]V, t Tuple[K, V]) map[K]V {
	if into == nil {
		into = make(map[K]V)
	}
	into[t.Key()] = t.Value()
	return into
}

// MakeSet is a FuncCollect that can take a Seq of comparable values and store them in a standard Go set (map[]struct{}).
// This function works with nil or a pre-built map[K]struct{} as initial value.
func MakeSet[K comparable](into map[K]struct{}, k K) map[K]struct{} {
	if into == nil {
		into = make(map[K]struct{})
	}
	into[k] = struct{}{}
	return into
}

// GroupBy is a FuncCollect that can take a Seq of Tuple values and group them by Tuple.Key in a map.
// All values of Tuple.Value are appended to a slice under each key.
// This function works with nil or a pre-built map[K][]V as initial value.
//
// Example, grouping serial numbers under a slice of first names:
//
//	names := SliceOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
//	tups := ZipOf[string, int](names, RangeFrom(0))
//	result := Reduce(nil, GroupBy[string, int], tups)
//
// Then the result is
//
//	map[string][]int{
//	  "bob":    {0, 2, 4},
//	  "alan":   {1, 5},
//	  "scotty": {3},
//	}
func GroupBy[K comparable, V any](into map[K][]V, tup Tuple[K, V]) map[K][]V {
	if into == nil {
		into = make(map[K][]V)
	}
	into[tup.Key()] = append(into[tup.Key()], tup.Value())
	return into
}

// UpdateMap is used to build a new FuncCollect that can update a map[K]V in place.
// It updates the element at Tuple.Key() with the provided FuncUpdate.
// Classic update functions could be fnmath.Max, fnmath.Min, or fnmath.Sum.
//
// Example, counting the number of unique names in a slice:
//
//	names := fn.SliceOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
//	tups := fn.ZipOf[string, int](names, Constant(1))
//	res := fn.Reduce(nil, fn.UpdateMap[string, int](fnmath.Sum[int]), tups)
//	fmt.Println(res)
//
// Prints:
//
//	map[alan:2 bob:3 scotty:1]
func UpdateMap[K comparable, V any](updater FuncUpdate[V]) FuncCollect[map[K]V, Tuple[K, V]] {
	return func(into map[K]V, tup Tuple[K, V]) map[K]V {
		if into == nil {
			into = make(map[K]V)
		}

		// Do the update
		into[tup.Key()] = updater(into[tup.Key()], tup.Value())
		return into
	}
}

// UpdateSlice is used to build a new FuncCollect that can update a slice []V in place.
// It looks at the elements in the index specified by Tuple.X(), ensures that the slice is big enough
// (growing it if needed), and updates the value at that index with the provided FuncUpdate.
// Classic update functions could be fnmath.Max, fnmath.Min, or fnmath.Sum.
func UpdateSlice[I constraints.Integer, V any](updater FuncUpdate[V]) FuncCollect[[]V, Tuple[I, V]] {
	return func(into []V, tup Tuple[I, V]) []V {
		idx := int(tup.Key())

		// Ensure target slice has required size (len(into) >= idx+1)
		if into == nil {
			into = make([]V, idx+1, 3*(idx+1)/2)
		} else if len(into) <= idx {
			if cap(into) >= idx {
				into = into[:idx+1] // into is big enough, we can just extend
			} else {
				// grow into
				newInto := make([]V, idx+1, 3*(idx+1)/2)
				copy(newInto, into)
				into = newInto
			}
		}

		// Do the update
		into[idx] = updater(into[idx], tup.Value())
		return into
	}
}

// OrderAsc is a FuncLess that can be used with Slice.Sort
func OrderAsc[T constraints.Ordered](t1, t2 T) bool {
	return t1 < t2
}

// OrderDesc is a FuncLess that can be used with Slice.Sort
func OrderDesc[T constraints.Ordered](t1, t2 T) bool {
	return t1 > t2
}

// OrderTupleAsc is a FuncLess that can be used with Slice.Sort
func OrderTupleAsc[K constraints.Ordered, V any](t1, t2 Tuple[K, V]) bool {
	return t1.Key() < t2.Key()
}

// OrderTupleDesc is a FuncLess that can be used with Slice.Sort
func OrderTupleDesc[K constraints.Ordered, V any](t1, t2 Tuple[K, V]) bool {
	return t1.Key() > t2.Key()
}

// IsZero is a Predicate that returns true if the input is the zero value of the type T.
// Can be used with Seq methods like Seq.TakeWhile, Seq.Where, and Seq.While.
func IsZero[T comparable](t T) bool {
	var zero T
	return t == zero
}

// IsNonZero is a Predicate that returns true if the input is a non-zero value of the type T.
// Can be used with Seq methods like Seq.TakeWhile, Seq.Where, and Seq.While.
func IsNonZero[T comparable](t T) bool {
	var zero T
	return t != zero
}

// Is returns a Predicate checking equality against the given argument.
// If you are checking for the zero value of a type it is faster to use IsZero.
func Is[T comparable](val T) Predicate[T] {
	return func(other T) bool {
		return val == other
	}
}

// IsNot returns a Predicate checking inequality against the given argument.
// If you are checking against the zero value of a type it is faster to use IsNonZero.
func IsNot[T comparable](val T) Predicate[T] {
	return func(other T) bool {
		return val != other
	}
}

// GreaterThanZero is a Predicate
func GreaterThanZero[T constraints.Ordered](t T) bool {
	var zero T
	return t > zero
}

// LessThanZero is a Predicate
func LessThanZero[T constraints.Ordered](t T) bool {
	var zero T
	return t < zero
}

// GreaterThan returns a Predicate. If you compare against zero then GreaterThanZero is more efficient.
func GreaterThan[T constraints.Ordered](val T) Predicate[T] {
	return func(other T) bool {
		return other > val
	}
}

// LessThan returns a Predicate. If you compare against zero then LessThanZero is more efficient.
func LessThan[T constraints.Ordered](val T) Predicate[T] {
	return func(other T) bool {
		return other < val
	}
}

// Not takes a Predicate and returns another predicate that is the logical inverse.
func Not[T any](pred Predicate[T]) Predicate[T] {
	return func(t T) bool {
		return !pred(t)
	}
}

// TupleWithKey creates a FuncMap to use with MappingOf or Map.
// The returned function yields Tuples keyed by the keySelector's return value.
// Usually used in conjunction with Reduce and MakeMap to build a map[K]V.
func TupleWithKey[K comparable, V any](keySelector FuncMap[V, K]) func(V) Tuple[K, V] {
	return func(v V) Tuple[K, V] {
		return TupleOf(keySelector(v), v)
	}
}

// Reduce executes a Seq, collecting the results via a collection function (FuncCollect).
// The method signature follows append() and copy() conventions,
// having the destination to put data into first.
// In other languages and libraries this function is also known as "fold".
//
// This library ships with a suite of standard collector functions.
// These include Append, MakeMap, MakeSet, MakeString, MakeBytes, Count,
// GroupBy, UpdateMap, UpdateSlice, fnmath.Sum, fnmath.Min, fnmath.Max,.
//
// The second argument, "into", can often be left as nil. It is the initial state for the collector.
// If you want to pre-allocate or reuse a buffer you can pass it in here. Or if you want to have
// a certain prefix on a string you can pass in a strings.Builder where you have added the prefix.
//
// If the seq produces and error the returned Opt will capture it, similarly if the seq is empty
// the returned opt.Opt will be empty.
func Reduce[T, E any](collector FuncCollect[T, E], into T, seq Seq[E]) opt.Opt[T] {
	empty := true
	tail := seq.ForEach(func(elem E) {
		into = collector(into, elem)
		empty = false
	})

	if err := Error(tail); err != nil {
		return opt.ErrorOf[T](err)
	} else if empty {
		return opt.Empty[T]()
	}

	return opt.Of(into)
}

// Do executes a Seq. The main use case is when you are primarily interested in triggering side effects.
// For parallel execution of Seqs please look at Go.
// In all normal applications the returned Seq will be empty. If ht e Seq is doing IO or other things
// with possibilities of runtime failures you may need to check it for errors with the Error function.
func Do[T any](seq Seq[T]) Seq[T] {
	return seq.ForEach(func(_ T) {})
}

// Any executes the Seq up to a point where the predicate returns true.
// If it finds such an element it returns true, otherwise if there are no matches, false.
// An empty seq will always return false.
//
// To check if a seq contains a zero element:
//
//		sq := seq.SliceOfArgs(-1, 0, 1, 2)
//		seq.Any(sq, seq.IsZero[int]) // returns true
//
//	 sq = SliceOf[int](nil)
//	 seq.Any(sq, seq.IsZero[int]) // returns false since sq is empty
func Any[T any](seq Seq[T], pred Predicate[T]) bool {
	fst, _ := seq.Where(pred).First()
	return fst.Ok()
}

// All executes the Seq and returns true iff all elements return true under the predicate.
// An empty seq will always return true.
//
// To check if all elements in a seq are non-zero:
//
//	 sq := seq.SliceOfArgs(1, 2, 3)
//		seq.All(sq, seq.IsNonZero[int]) // returns true
//
//	 sq = seq.SliceOf[int](nil)
//		seq.All(sq, seq.IsNonZero[int]) // returns true since sq is empty
func All[T any](seq Seq[T], pred Predicate[T]) bool {
	fstMismatch, _ := seq.Where(Not(pred)).First()
	return fstMismatch.Empty()
}

// Last executes the Seq and returns the last element or an empty Opt.
func Last[T any](seq Seq[T]) opt.Opt[T] {
	var last T
	var i int
	tail := seq.ForEach(func(t T) {
		last = t
		i++
	})
	if i > 0 {
		return opt.Of(last)
	}
	if err := Error(tail); err != nil {
		return opt.ErrorOf[T](err)
	}
	return opt.Empty[T]()
}

// ErrNotOne is returned from One if the seq contains more than 1 element.
var ErrNotOne = errors.New("sequence contains more than 1 element")

// One returns a valid opt if the seq contains exactly one element.
// If the seq is empty the option will be invalid with opt.ErrEmpty as usual,
// and if there is more than 1 element then the error will be ErrNotOne.
func One[T any](seq Seq[T]) opt.Opt[T] {
	fst, tail := seq.First()
	if fst.Empty() {
		return fst
	}
	second, _ := tail.First()
	if second.Ok() {
		return opt.ErrorOf[T](ErrNotOne)
	}
	return fst
}

// IsEmpty returns true if the seq is empty.
// This function may execute the first element of the seq, if the length can not be determined.
//
// Checking if a seq is empty is rarely necessary. All operations are valid on an empty sequence,
// even nil is a valid Slice. Just consume the seq and check the resulting opt.Opt or Slice.
// Note than you can use the normal len() function on Slice, Set, Map, Chan, and String.
func IsEmpty[T any](seq Seq[T]) bool {
	if sz, ok := seq.Len(); ok {
		return sz == 0
	} else if sz == LenInfinite {
		return false
	}

	// No well-defined length we need to execute the first element
	fst, _ := seq.First()
	return fst.Empty()
}
