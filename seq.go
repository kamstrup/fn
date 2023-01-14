package fn

// LenUnknown is returned by Seq.Len() for Seqs that do not have a well-defined length.
// A typical example will be a Seq returned by Seq.Where(). Since the Predicate of
// the resultant Seq is lazily executed the length can not be known up front.
const LenUnknown = -1

// LenInfinite is returned by certain Seqs such as SourceOf(), when the sequence will
// not terminate by itself but continue to yield values indefinitely.
const LenInfinite = -2

// Seq is the primary interface for the fn() library.
// Seqs should be thought of a lazily computed collections of elements.
// Operations that force traversing or computation of the Seq are said to "execute" the Seq.
// As a rule of thumb any method that returns an Array will execute that part of the Seq.
// For example, seq.Take(7) executes the first 7 elements and returns them in an Array,
// and the rest of the Seq is untouched and returned as a tail Seq.
// Any method that executes the Seq must document it explicitly.
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
	Len() (int, bool)
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
	// To retrieve the last element in a Seq you can use Last.
	First() (Opt[T], Seq[T])
	// Map lazily converts elements of the Seq into a value of the same type.
	// Classic examples would be to convert strings to lowercase, multiply a range of numbers by Pi, and similar.
	// If you need to change the type of the elements you must use the function fn.MapOf(),
	// since Go does not support type parameters on methods.
	Map(funcMap FuncMap[T, T]) Seq[T]
}

func errOrEmpty[T any](seq Seq[T]) Seq[T] {
	if err := Error(seq); err != nil {
		return ErrorOf[T](err)
	}
	return SeqEmpty[T]()
}
