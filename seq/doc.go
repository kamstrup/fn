// Package seq provides functions and data structures for incorporating
// functional programming ideas into standard Go programs.
//
// The key data structure is the Seq, which is short for "sequence".
// Seqs are generally lazy and stateless, unless explicitly noted otherwise.
// All methods return a new Seq that you must work on, similar to how append() works.
//
// There are seqs for working with slices (see SliceOf, SliceOfArgs),
// maps and sets (see MapOf and SetOf), channels (see ChanOf), numeric ranges
// (see RangeOf), or created from generator functions (see SourceOf).
// Seqs can be joined with ConcatOf and FlattenOf, split with SplitOf,
// iterated in parallel as Tuple with ZipOf.
package seq
