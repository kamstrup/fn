package fn

import (
	"bytes"
	"strings"

	"github.com/kamstrup/fn/constraints"
)

type Func0 func()

type Func0Err func() error

type Func1[T any] func(T)

type Func1Err[T any] func(T) error

type Func2[S, T any] func(S, T)

type Func2Err[S, T any] func(S, T) error

// FuncMap is a function mapping type S to type T.
// Used with fx. MapOf(), OptMap(), and TupleWithKey().
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

// Sum is a FuncCollect and a FuncUpdate, for use with Into or UpdateAt, that sums up the elements it sees.
func Sum[T constraints.Arithmetic](into, t T) T {
	return into + t
}

// Max is a FuncCollect and a FuncUpdate, for use with Into or UpdateAssoc, that returns the maximal element.
func Max[T constraints.Ordered](s, t T) T {
	if s > t {
		return s
	}
	return t
}

// Min is a FuncCollect and a FuncUpdate, for use with Into or UpdateAssoc, that returns the minimal element.
func Min[T constraints.Ordered](s, t T) T {
	if s < t {
		return s
	}
	return t
}

// Count is a FuncCollect for use with Into, that counts the number of elements it sees.
func Count[T any](into int, _ T) int {
	return into + 1
}

// Append is a FuncCollect for use with Into, that uses the standard Go append() function.
// This function works with nil or a pre-built slice as initial value.
func Append[T any](into []T, t T) []T {
	return append(into, t)
}

// MakeString is a FuncCollect for use with Into that writes strings into a bytes.Buffer.
// This function works with nil or a pre-built bytes.Buffer as initial value.
func MakeString(into *strings.Builder, s string) *strings.Builder {
	if into == nil {
		into = &strings.Builder{}
	}
	_, _ = into.WriteString(s)
	return into
}

// MakeBytes is a FuncCollect for use with Into that writes bytes into a bytes.Buffer.
// This function works with nil or a pre-built bytes.Buffer as initial value.
func MakeBytes(into *bytes.Buffer, b []byte) *bytes.Buffer {
	if into == nil {
		into = &bytes.Buffer{}
	}
	_, _ = into.Write(b)
	return into
}

// MakeAssoc is a FuncCollect that can take a Seq of Tuple elements and store them in a standard Go map.
// This function works with nil or a pre-built map[K]V as initial value.
func MakeAssoc[K comparable, V any](into map[K]V, t Tuple[K, V]) map[K]V { // FIXME MAYBE USE ASSOC instead of map
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
//	names := ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
//	tups := ZipOf[string, int](names, NumbersFrom(0))
//	result := Into(nil, GroupBy[string, int], tups)
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

// UpdateAssoc is used to build a new FuncCollect that can update a map[K]V in place.
// It updates the element at Tuple.Key() with the provided FuncUpdate.
// Classic update functions could be Max, Min, or Sum.
//
// Example, counting the number of unique names in a slice:
//
//	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
//	tups := fn.ZipOf[string, int](names, Constant(1))
//	res := fn.Into(nil, fn.UpdateAssoc[string, int](Sum[int]), tups)
//	fmt.Println(res)
//
// Prints:
//
//	map[alan:2 bob:3 scotty:1]
func UpdateAssoc[K comparable, V any](updater FuncUpdate[V]) FuncCollect[map[K]V, Tuple[K, V]] {
	return func(into map[K]V, tup Tuple[K, V]) map[K]V {
		if into == nil {
			into = make(map[K]V)
		}

		// Do the update
		into[tup.Key()] = updater(into[tup.Key()], tup.Value())
		return into
	}
}

// UpdateArray is used to build a new FuncCollect that can update a slice []V in place.
// It looks at the elements in the index specified by Tuple.X(), ensures that the slice is big enough
// (growing it if needed), and updates the value at that index with the provided FuncUpdate.
// Classic update functions could be Max, Min, or Sum.
func UpdateArray[I constraints.Integer, V any](updater FuncUpdate[V]) FuncCollect[[]V, Tuple[I, V]] {
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

// OrderAsc is a FuncLess that can be used with Array.Sort
func OrderAsc[T constraints.Ordered](t1, t2 T) bool {
	return t1 < t2
}

// OrderDesc is a FuncLess that can be used with Array.Sort
func OrderDesc[T constraints.Ordered](t1, t2 T) bool {
	return t1 > t2
}

// OrderTupleAsc is a FuncLess that can be used with Array.Sort
func OrderTupleAsc[K constraints.Ordered, V any](t1, t2 Tuple[K, V]) bool {
	return t1.Key() < t2.Key()
}

// OrderTupleDesc is a FuncLess that can be used with Array.Sort
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

// TupleWithKey creates a FuncMap to use with MapOf or OptMap.
// The returned function yields Tuples keyed by the keySelector's return value.
// Usually used in conjunction with Into and MakeAssoc to build a map[K]V.
func TupleWithKey[K comparable, V any](keySelector FuncMap[V, K]) func(V) Tuple[K, V] {
	return func(v V) Tuple[K, V] {
		return TupleOf(keySelector(v), v)
	}
}

// Into executes a Seq, collecting the results via a collection function (FuncCollect).
// The method signature follows append() and copy() conventions,
// having the destination to put data into first.
// In other languages and libraries this function is also known as "reduce" or "fold".
//
// This library ships with a suite of standard collector functions.
// These include Append, MakeAssoc, MakeSet, MakeString, MakeBytes, Sum, Count, Min, Max,
// GroupBy, UpdateAssoc, UpdateArray.
//
// The first argument, "into", can often be left as nil. It is the initial state for the collector.
// If you want to pre-allocate or reuse a buffer you can pass it in here. Or if you want to have
// a certain prefix on a string you can pass in a strings.Builder where you have added the prefix.
//
// If the seq produces and error the returned Opt will capture it, similarly if the seq is empty
// the returned Opt will be empty.
func Into[T, E any](into T, collector FuncCollect[T, E], seq Seq[E]) Opt[T] {
	empty := true
	tail := seq.ForEach(func(elem E) {
		into = collector(into, elem)
		empty = false
	})

	if err := Error(tail); err != nil {
		return OptErr[T](err)
	} else if empty {
		return OptEmpty[T]()
	}

	return OptOf(into)
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
func Any[T any](seq Seq[T], pred Predicate[T]) bool {
	fst, _ := seq.Where(pred).First()
	return fst.Ok()
}

// All executes the Seq and returns true iff all elements return true under the predicate.
func All[T any](seq Seq[T], pred Predicate[T]) bool {
	fstMismatch, _ := seq.Where(Not(pred)).First()
	return fstMismatch.Empty()
}

// Last executes the Seq and returns the last element or an empty Opt.
func Last[T any](seq Seq[T]) Opt[T] {
	var last T
	var i int
	tail := seq.ForEach(func(t T) {
		last = t
		i++
	})
	if i > 0 {
		return OptOf(last)
	}
	if err := Error(tail); err != nil {
		return OptErr[T](err)
	}
	return OptEmpty[T]()
}
