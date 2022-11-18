Fn(), Functional Programming for Golang
====
Fn is library for golang that enables functional programming techniques
to be blended it with standard idiomatic Go code.

Philosophy
----
 * Be pragmatic. Fn is not idealistic, academic, or purely functional.
   The aim is blend cleanly in to existing Go programs, not to invent
   a new paradigm for Go. While most of Fn is built on immutable functional
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
// Zip(s1 Seq[S], s2 Seq[T]) Seq[Tuple[S,T]]
// Concat(Seq[T], Seq[T]) Seq[T]
// seq.Any(pred)/All(pred)
// seq.Split(pred) Seq[Seq[T]]
// Error handling: fn.Must()
// Seq over a channel
// Seq over a Set, and a Collector func for creating a Set
// Select on channel
// Source(FuncSource), Range(int, int), Tuple[S,T] as Seq?
// Reverse?
// testing utils assert/require?
// seq.Go(n, f) (n goroutines). Auto-wait, or SeqGo.Wait()?
```

DONE
---
```
Seqs of slices (Array) and maps (Assoc).
Into() and various Collector funcs for it to aggregate Seqs into numbers or new collections
Tuples for help with building maps (and Zip in the future)
Lazily evaluated Where() and Map()
Opt with Try() and Must()
Some simple helpers to write testing.T tests using Seqs.
``` 