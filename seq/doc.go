// Package fn provides functions and data structures for incorporating
// functional programming ideas into standard Go programs.
//
// The key data structure is the Seq, which is short for "sequence".
// Seqs are generally lazy and stateless. All methods return a new Seq
// that you must work on, similar to how append() works.
//
// Another important interface is the Opt which is short for "option".
// Opts represent a value that may or may not be there. They may also
// encapsulate an error.
//
// There are seqs for working with slices (see SliceOf, SliceOfArgs),
// maps and sets (see AssocOf and SetOf), channels (see ChanOf), numeric ranges
// (see RangeOf), or created from generator functions (see SourceOf).
// Seqs can be joined with ConcatOf and FlattenOf, split with SplitOf,
// iterated in parallel as Tuple with ZipOf.
//
// There are sub-packages that expands what seqs can be applied to, like fnmath and fnio.
// Finally, there is a light limited-scope function library called fx.
package seq
