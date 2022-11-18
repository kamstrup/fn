Fn(), Functional Programming for Golang
====
Fn is library for golang that enables functional programming techniques
to be blended it with standard idiomatic Go code.

Philosophy
----
 * Be pragmatic. Fn() is not idealistic, academic, or purely functional.
   The aim is blend cleanly in to existing Go programs, not to invent
   a new paradigm for Go. While most of Fn() is built on immutable functional
   data structures, we live in the mutable world of Go. We try to blend in
   smoothly instead of ice skating uphill.
 * Be small and concise. Fn will not come with a tonne of data structures
   that you might expect from a full-featured functional library. The
   core Fn library does not come with tree structures or persistent hash tables, fx.
 * No dependencies. Only the Go standard library.
 * Include simple affordances to interop with constructs from the Go standard
   library, but no big new frameworking for doing IO or other stuff.

TODO
---
```
// ConcatOf(seqs Seq[Seq[T]]) Seq[T], and ConcatOfArgs(seqs ... Seq[T]) Seq[T]
// seq.Any(pred)/All(pred)
// seq.Split(pred) Seq[Seq[T]]
// Seq over a channel
// Select on channel
// Range(int, int)
// Tuple[S,T] as Seq[any]? (we have to do "any" bc the types S!=T)
// Improve testing utils assert/require? Move to own package?
// seq.Go(n, f) (n goroutines) and seq.GoErr(, f). Auto-wait, or SeqGo.Wait()? Control chan? 
// MergeSort[T any](FuncLess[T], seqs ... Seq[T]) Seq[T] -- lazy merge sorting of pre-sorted Seqs
// Seq.Last()
// Better testing for Zip
```

DONE
---
```
Seqs of slices (Array) and maps (Assoc and Set).
Into() and various Collector funcs for it to aggregate Seqs into numbers or new collections
Tuples for help with building maps (and Zip in the future)
Lazily evaluated Where() and Map()
Opt with Try() and Must()
Some simple helpers to write testing.T tests using Seqs.
Sorting with some pre-declared generic helpers
Infinite Seqs from source functions
Zip 2 Seqs into a Seq of Tuples
```