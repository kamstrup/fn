About Fn
====
Fn is library for golang that enable you to blend functional programming techniques
with standard idiomatic Go code.

We are inspired by [Clojure](https://github.com/clojure/clojure),
[Vavr](https://github.com/vavr-io/vavr), and the Java Streams APIs that were
introduced back in Java 8, and want to provide something of similar spirit
that makes it even more fun to write Go code.

Philosophy
====
* Be pragmatic. Fn() is not idealistic, academic, or purely functional.
  The aim is blend cleanly into existing Go programs, not to invent
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

[1]: Some of the more advanced functional data types, like the ones mentioned,
are definitely super useful and would fit well in some extension library for Fn().
(or perhaps a sub-package, let's see <3)

## Terminology
For starters let's get some terminology in place.

[Seq](seq.md): The core data structure is `Seq[T]`. It is short for Sequence. It can encapsulate any kind of data structure
where elements can be traversed sequentially (but not necessarily ordered).

[Slice](slice.md): Standard Go slices `[]T` are wrapped in the `seq.Slice[T]` type.
The `Slice` type is a public subtype of `[]T` so you can do numeric indexing on a `Slice`,
and use `cap()`, `len()`, and for-range loops on them.
Slices are seqs, but also add some extra methods like `Sort()` and `Reverse()`.

**Tuple:** Or "pair". Represents two data points. A helper mainly used when working with maps,
where the tuple captures a key and a value. Maps can be interpreted as a seq of tuples, or if you
build a seq of tuples you can create map from it. Tuples also show up when "zipping" 2 seqs.

[Opt](opt.md): Returned from operations where you are not certain to get a result.
For example when you call `sq.First()`. If the seq is empty you get back an empty opt,
and an empty tail sq. Opts can also be used to encapsulate errors.

**sq.Len():** Lengths are handled in a special way in Fn(). They are allowed to be finite,
unknown, or infinite. Making these distinctions opens the possibility of pre-allocating
slices and maps of the correct size, which can make a big difference in performance critical code.
In most circumstances you will not need to use Len(). All operations are valid on empty seqs and
empty opts, so just set up your pipeline of operations and check if the end result is valid.
The `Len()` method is rarely needed. The types `seq.Map`, `seq.Set`, `seq.Slice`, and `seq.Chan` all
support the builtin `len()` function.


## API Overview
If you just want to jump in and see some code you can check out
[the simple examples](https://github.com/kamstrup/fn/blob/main/examples/simple_test.go).
Otherwise, here follows a brief overview.

Fn also bundles a [very simple sub-library called 'slice'](#the-slice-package)
that you can use to do 1-line functional constructs.


