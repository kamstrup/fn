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
   