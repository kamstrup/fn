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

TODO
---
```
DOCS
* Put examples in this README

API CHANGES
* Error reporting for fn.Seq.Array()? Tricky since Array is not a struct, but a just slice type alias
* Rename fn.Set to Uniq? ("set" is a heavily overloaded word)
* Rename Seq.Shape() to .Map()

FEATURES (in order of prio)
* WIP A small IO package "fnio" to help walking an io.Reader as a Seq[[]byte], and same for writing?
* fnio.DirOf(dirName), * fnio.DirTreeOf(dirName) (recursive)
* fnio.LinesOf(io.Reader)
* seq.Limit(n) Seq[T], lazy counterpart to seq.Take(n)
* Compare func helpers LessThan, GreaterThan, Is, IsNot
* seq.Split(pred) Seq[Seq[T]] (including or excluding the separator? We need both modes)
* RunesOf(string) Seq[rune]
* Improve testing utils assert/require? Move to own package?
* A small JSON package "fnjson" to help reading and writing Seqs of JSON objects
* Array.Reverse() Seq[T], zero-copy reversed Array (just reversed view on arr)
* MultiChan() Seq that selects on multiple chan T?
* fn.GroupBy(Seq[S], FuncMap[S,T]) map[T][]S
* Something for context.Context? Support cancel() cb and Done() chans? fncontext package...
* fn.GoErr(seq, numTasks, FuncMapErr) -- or some version of fn.Go() with cancellation and error handling. 
* Seq.Last() maybe? 
* Tuple[S,T] as Seq[any]? (we have to do "any" bc the types S!=T)
* Set, Assoc are just straight type wrappers. Make it public API like for Array and String?
* MergeSort[T any](FuncLess[T], seqs ... Seq[T]) Seq[T] -- lazy merge sorting of pre-sorted Seqs

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
```

Performance
----
If the foundational functional data structures and algorithms is not done carefully,
execution speed and memory usage will suffer. Fn() is designed to make the best of what
the Go runtime can provide. Initial benchmarks puts it as a top contender among Golang
functional libraries. See benchmarks here https://github.com/mariomac/go-stream-benchmarks/pull/1