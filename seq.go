package fn

// LenUnknown is returned by Seq.Len() for Seqs that do not have a well-defined length.
// A typical example will be a Seq returned by Seq.Where(). Since the Predicate of
// the resultant Seq is lazily executed the length can not be known up front.
const LenUnknown = -1

// Seq is the primary interface for the fn() library.
// Seqs should be thought of a lazily computed collections of elements.
// Operations that force traversing or computation of the Seq are said to "execute" the Seq.
// As a rule of thumb any method that returns an Array will execute that part of the Seq.
// For example, seq.Take(7) executes the first 7 elements and returns them in an array,
// and the rest of the Seq is untouched and returned as a tail Seq.
// Any method that executes the Seq must document it explicitly.
type Seq[T any] interface {
	// ForEach executes the Seq and calls f on each element.
	ForEach(f Func1[T])
	// ForEachIndex executes the Seq and calls f on each index and element.
	ForEachIndex(f Func2[int, T])
	// Len returns the number of elements in the Seq OR LenUnknown if the Seq does not have a well-defined length.
	Len() int
	// Array executes the Seq and stores all elements in memory as an Array
	Array() Array[T]
	// Take executes up to the first N elements of the Seq and returns the rest in a tail Seq
	Take(int) (Array[T], Seq[T])
	// TakeWhile executes the Seq while Predicate returns true,
	// then returns those elements in an Array and the rest in a tail Seq.
	// This library ships with a few in-built predicates, like fx, IsZero and IsNonZero.
	TakeWhile(predicate Predicate[T]) (Array[T], Seq[T])
	// Skip drops (up to) the first N elements in the Seq, executing them, and returns a tail Seq.
	Skip(int) Seq[T]
	// Where returns a Seq that only includes elements where Predicate returns true.
	// This library ships with a few in-built predicates, like fx, IsZero and IsNonZero.
	Where(Predicate[T]) Seq[T]
	// While returns a Seq with the initial series of elements where predicate returns true.
	// This library ships with a few in-built predicates, like fx, IsZero and IsNonZero.
	While(predicate Predicate[T]) Seq[T]
	// First executes the first element and returns an Opt with it. The tail is returned as a Seq.
	First() (Opt[T], Seq[T])
	// All executes the Seq and returns true iff all elements return true under the predicate.
	All(Predicate[T]) bool
	// Any executes the Seq up to a point where the predicate returns true.
	// If it finds such an element it returns true, otherwise if there are no matches, false.
	Any(Predicate[T]) bool
}

func seqAny[T any](seq Seq[T], pred Predicate[T]) bool {
	fst, _ := seq.Where(pred).First()
	return fst.Ok()
}

func seqAll[T any](seq Seq[T], pred Predicate[T]) bool {
	fstMismatch, _ := seq.Where(Not(pred)).First()
	return fstMismatch.Empty()
}
