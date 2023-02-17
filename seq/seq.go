package seq

import "github.com/kamstrup/fn/opt"

// LenUnknown is returned by Seq.Len() for Seqs that do not have a well-defined length.
// A typical example will be a Seq returned by Seq.Where(). Since the Predicate of
// the resultant Seq is lazily executed the length can not be known up front.
const LenUnknown = -1

// LenInfinite is returned by certain Seqs such as SourceOf(), when the sequence will
// not terminate by itself but continue to yield values indefinitely.
const LenInfinite = -2

// Seq is the primary interface for the seq package.
// Seqs should be thought of as a stateless lazily computed collection of elements.
// Operations that force traversing or computation of the Seq are said to "execute" the Seq.
// As a rule of thumb any method that returns a Slice will execute that part of the Seq.
// For example, seq.Take(7) executes the first 7 elements and returns them in a Slice,
// and the rest of the Seq is untouched and returned as a tail Seq.
// Any method that executes the Seq must document it explicitly.
//
// Seq implementations that are not stateless or lazy must explicitly document that
// they are stateful and/or eager.
type Seq[T any] interface {
	// ForEach executes the Seq and calls f on each element.
	// Returns an empty Seq. If the Seq has capabilities for errors,
	// the returned Seq should be checked with Error().
	ForEach(f Func1[T]) Seq[T]
	// ForEachIndex executes the Seq and calls f on each index and element.
	// Returns an empty Seq. If the Seq has capabilities for errors,
	// the returned Seq should be checked with Error().
	ForEachIndex(f Func2[int, T]) Seq[T]
	// Len returns the number of elements in the Seq if it is well-defined.
	// If the boolean return value is false, the length is not well-defined and is either
	// LenUnknown or LenInfinite.
	//
	// Len must always be O(1) and must not execute the seq.
	//
	// You rarely want to check the length of a Seq. It is mostly intended for cases where
	// we can pre-allocate suitable size buffers to hold results. If you are looking to check
	// if a Seq contains a given element, or satisfies some property, you can use Any, All,
	// Last, One, or IsEmpty.
	Len() (int, bool)
	// ToSlice executes the Seq and stores all elements in memory as a Slice.
	// Recall that all functions and operations that works on a normal slice []T, also work directly on a Slice[T].
	ToSlice() Slice[T]
	// Limit returns a lazy Seq with maximally n elements.
	// The Take method is related, but Take is different in that it executes the first n elements and returns the tail.
	// Limit is lazy, and can as such not return a tail.
	Limit(n int) Seq[T]
	// Take executes up to the first N elements of the Seq and returns the rest in a tail Seq.
	// The Limit method is related, but different because Limit is lazy and does not return a tail.
	Take(int) (Slice[T], Seq[T])
	// TakeWhile executes the Seq while Predicate returns true,
	// then returns those elements in a Slice and the rest in a tail Seq.
	// This library ships with a few in-built predicates, like fx, IsZero and IsNonZero.
	TakeWhile(predicate Predicate[T]) (Slice[T], Seq[T])
	// Skip drops (up to) the first N elements in the Seq, executing them, and returns a tail Seq.
	Skip(int) Seq[T]
	// Where returns a Seq that only includes elements where Predicate returns true.
	// This library ships with a few in-built predicates, like fx, IsZero and IsNonZero.
	Where(Predicate[T]) Seq[T]
	// While returns a Seq with the initial series of elements where predicate returns true.
	// This library ships with a few in-built predicates, like fx, IsZero and IsNonZero.
	While(predicate Predicate[T]) Seq[T]
	// First executes the first element and returns an Opt with it. The tail is returned as a Seq.
	// To retrieve the last element in a Seq you can use Last.
	First() (opt.Opt[T], Seq[T])
	// Map lazily converts elements of the Seq into a value of the same type.
	// Classic examples would be to convert strings to lowercase, multiply a range of numbers by Pi, and similar.
	// If you need to change the type of the elements you must use the function fn.MappingOf(),
	// since Go does not support type parameters on methods.
	Map(funcMap FuncMap[T, T]) Seq[T]
}

func errOrEmpty[T any](seq Seq[T]) Seq[T] {
	if err := Error(seq); err != nil {
		return ErrorOf[T](err)
	}
	return Empty[T]()
}
