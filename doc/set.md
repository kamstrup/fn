Working with Sets
====
One of the most important things to realize about sets in Fn is that `seq.Set`
is just a type wrapper for `map[T]struct{}`. This means that standard Go builtins and syntax
work on it:

 * `make(seq.Set[string], 10)`
 * Literal assignment `mySet := seq.Set[string]{"one": {}, "two": {}}`
 * `len(mySet)`
 * for-range loops
 * Indexing, like `_, hasElement := mySet[elem]`

Limitations and Caveats
----
At the time of release of Fn 1.0 the Go compiler's support for generics is very
limited. This means that it can not tell that a `seq.Set[T]` implements the `seq.Seq[T]`
interface. Either construct your sets with `seq.SetOf()` which returns a `seq.Seq[T]` or
call `mySet.Seq()` to explicitly cast it.

Creating Sets
----
A very short summary of different ways to create a set:

 * `seq.SetAs()`
 * `seq.SetAsArgs(...)`
 * `seq.Set[string]{"one": {}, "two": {}}`
 * `nil` is a valid empty set
 * `make(seq.Set[string], 10)`
 * `seq.Reduce()` with the `seq.MakeSet` collection function

Set Features
----
Apart from being usable as a standard Go map, and a seq, sets have a few extra methods:

 * Set.Union() returns a lazy union of 2 sets
 * Set.Intersect() returns a lazy intersection of 2 sets.
 * Set.Contains(k)