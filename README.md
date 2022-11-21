Fn(), Functional Programming for Golang
====
Fn is library for golang that enables functional programming techniques
to be blended it with standard idiomatic Go code.

We are inspired by Clojure and the Java Streams APIs that were
introduced back in Java 8, and want to provide something of similar spirit,
that makes it even more fun to write GO code.

Philosophy
----
 * Be pragmatic. Fn() is not idealistic, academic, or purely functional.
   The aim is blend cleanly in to existing Go programs, not to invent
   a new paradigm for Go. While most of Fn() is built on immutable functional
   data structures, we live in the mutable world of Go. We try to blend in
   smoothly instead of ice skating uphill.
 * Be small and concise. Fn() will not come with a tonne of data structures
   that you might expect from a full-featured functional library. The
   core Fn() library does not come with immutable red-black tree structures or
   persistent hash tables, fx.[1]
 * No dependencies. Only the Go standard library.
 * Include simple affordances to interop with constructs from the Go standard
   library, but no big new frameworking for doing IO or other stuff.
 * If we start talking about "monoids", "transducers", and even simpler terms
   like "fold", and "reduce", 75% of developers will start to zone out or just walk away.
   We prioritize simple, well-known, terms that produce readable code that most
   developers will have a chance of understanding and enjoying.

[1]: Some of the more advanced functional data types, like the ones mentioned,
are definitely super useful and would fit well in some extension library for Fn().
(or perhaps a sub-package, let's see <3)

TODO
---
```
// STARTED ConcatOf(seqs Seq[Seq[T]]) Seq[T], and ConcatOfArgs(seqs ... Seq[T]) Seq[T]
// seq.Limit(n) Seq[T], lazy counterpart to seq.Take(n)
// Compare func helpers LessThan, GreaterThan, Is, IsNot
// seq.Any(pred)/All(pred) (should any return a tail Seq, or just bool)?
// seq.Split(pred) Seq[Seq[T]]
// Seq over a channel (maybe some special prupose helpers for Select on channel?)
// Do we need a special string Seq? StringOf(string)
// Tuple[S,T] as Seq[any]? (we have to do "any" bc the types S!=T)
// Improve testing utils assert/require? Move to own package?
// seq.Go(n, f) (n goroutines) and seq.GoErr(, f). Auto-wait, or SeqGo.Wait()? Control chan? 
// MergeSort[T any](FuncLess[T], seqs ... Seq[T]) Seq[T] -- lazy merge sorting of pre-sorted Seqs
// Seq.Last()
// Better testing for Zip
// Seq of single element
// EmptySeq impl. (currently just wraps an empty slice), but an empty struct{} would do even better
// Maybe a "Random Access"[K,V] interface that Array, AssocOf, and SetOf can implement
//         (although Array is just a []T, Assoc just a map[K]V, and Set a map[K]struct{},
            so support random access via subscripts directly.) 
// Should Len() return (int, ok) instead? DO we want a special LenInfinite?
           (ZipOf() would benefit in pre-allocs by knowing if one of the seqs were infinite)
```

DONE
---
```
Seqs of slices (Array) and maps (Assoc and Set).
Into() and various Collector funcs for it to aggregate Seqs into numbers or new collections
Tuples for help with building maps
Lazily evaluated Where(), While(), and MapOf()
Opt with Try() and Must()
Some simple helpers to write testing.T tests using Seqs.
Sorting with some pre-declared generic helpers
Infinite Seqs from source functions
Zip 2 Seqs into a Seq of Tuples
RangeOf(from, to) in both directions
```