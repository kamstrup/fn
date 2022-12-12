Fn(), Functional Programming for Golang
====
Fn is library for golang that enables functional programming techniques
to be blended it with standard idiomatic Go code.

We are inspired by Clojure and the Java Streams APIs that were
introduced back in Java 8, and want to provide something of similar spirit,
that makes it even more fun to write Go code.

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

Fn() Quick Start
---
For starters let's get some terminology in place.

**Seq:** The core data structure is `Seq[T]`. It is short for Sequence.
The Seq API is designed to work on top of immutable structures,
thus there is no stateful "iterator". Walking through a Seq is done similarly 
to how you `append()` elements to a slice in Go, but inversely.

```
ints = append(ints, i)
// OR if we call the slice "tail" and the element "head":
tail = append(tail, head)
```
The "inverse" of this operation looks like:
```
head, tail = tail.First() // pops the first element and returns a new tail
```
There are many easier ways to walk a Seq though. For example via `seq.ForEach()`
, `seq.Take()`, and `fn.Into()`

Generally seqs are immutable. Any exception to this will be clearly documented. 

[Check the interface definition for Seq here](https://github.com/kamstrup/fn/blob/main/seq.go). 

**Array:** Standard Go slices `[]T` are wrapped in the `fn.Array[T]` type.
The `Array` type is a public subtype of `[]T` so you can do numeric indexing on an `Array`.
Arrays are seqs, but also add some extra methods like `Sort()` and `Reverse()`.

**Assoc:** Because "map" is an overloaded word in functional programming,
Fn() uses the word "assoc" instead of `map[K]V` (pronounced with a soft "ch" at the end).
The word "map" is reserved for the mapping operation used to convert a Seq to something else.

**Creating a Seq**: A new Seq is obtained by one of the included constructor methods,
such as `ArrayOf()`, `ArrayOfArgs()`, `AssocOf()`, `MapOf()`, `ConcatOf()`, `ConcatOfArgs()`,
`RangeOf()`, `SetOf()`, `SourceOf()`, `ZipOf()`, and `StringOf()`.

**Tuple:** Or "pair". Represents to data points. A helper mainly used when working with assocs,
where the tuple captures a key and a value.

**Opt:** Returned from operations where you are not certain to get a result.
For example when you call `seq.First()`. If the seq is empty you get back an empty opt,
and an empty tail seq.

**Seq.Len():** Lengths are handled in a special way in Fn(). They are allowed to be finite, 
unknown, or infinite. 

TODO
---
```
DOCS
* Put examples in this README

API CHANGES
* Error reporting for fn.Seq.Array()? Tricky since Array is not a struct, but a just slice type alias
* ?? Array.Reverse() Seq[T], zero-copy reversed Array (just reversed view on arr). Currently modifies in place.

FEATURES (in order of prio)
* WIP A small IO package "fnio" to help walking an io.Reader as a Seq[[]byte], and same for writing?
* seq.Split(pred) Seq[Seq[T]] (including or excluding the separator? We need both modes)
* fnio.DirOf(dirName), * fnio.DirTreeOf(dirName) (recursive)
* fnio.LinesOf(io.Reader)
* seq.Limit(n) Seq[T], lazy counterpart to seq.Take(n)
* RunesOf(string) Seq[rune]
* A small JSON package "fnjson" to help reading and writing Seqs of JSON objects
* Special seqs for Assoc.Keys() and Assoc.Values()
* MultiChan() Seq that selects on multiple chan T?
* Something for context.Context? Support cancel() cb and Done() chans? fncontext package...
* fn.GoErr(seq, numTasks, FuncMapErr) -- or some version of fn.Go() with cancellation and error handling. 
* Seq.Last() maybe? 
* Tuple[S,T] as Seq[any]? (we have to do "any" bc the types S!=T)
* MergeSort[T any](FuncLess[T], seqs ... Seq[T]) Seq[T] -- lazy merge sorting of pre-sorted Seqs
* Compound FuncCollect, CollectorOf[S,T any](funcs ... FuncCollect[S,T]) FuncCollect[S,[]T]

OPTIMIZATIONS
* Seq of single element (see SingletOf(t))
* EmptySeq impl. (currently just wraps an empty slice), but an empty struct{} would do even better
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
Concat(Seq[Seq[T]]), and ConcatOfArgs(seqs ... Seq[T])
seq.Any() and seq.All()
Seq length is optional. sz, ok := seq.Len() and can be finite, infinite, or unknown 
Seq over a channel
Seq[byte] of a string
fn.Go() to execute a Seq in N goroutines and collect the results into another Seq
Special "error seq" that can be returned from IO ops and anything that can fail at runtime
RangeOf Seqs always have a well-defined length. Works for unsigned types as well
A bunch of helper predicates like Is(), IsNot(), IsZero(), IsGreaterThan(), etc
```

Performance
----
If the foundational functional data structures and algorithms is not done carefully,
execution speed and memory usage will suffer. Fn() is designed to make the best of what
the Go runtime can provide. Initial benchmarks puts it as a top contender among Golang
functional libraries. See benchmarks here https://github.com/mariomac/go-stream-benchmarks/pull/1