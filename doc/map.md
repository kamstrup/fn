Working with Maps
====
One of the most important things to realize about maps in Fn is that `seq.Map[K,V]`
is just a type wrapper for `map[K]V`. This means that standard Go builtins
work on it:

* `make(seq.Map[string,int], 10)`
* Literal assignment `myMap := seq.Map[string,int]{"one": 1, "two": 2}`
* `len(myMap)`
* for-range loops
* Indexing, like `val, hasElement := mySet[key]`

Limitations and Caveats
----
At the time of release of Fn 1.0 the Go compiler's support for generics is very
limited. This means that it can not tell that a `seq.Map[K,V]` implements the `seq.Seq[Tuple[K,V]]`
interface. Either construct your maps with `seq.MapOf()` which returns a `seq.Seq[Tuple[K,V]]` or
call `myMap.Seq()` to explicitly cast it.

Creating Maps
----
A very short summary of different ways to create a set:

* `seq.MapAs()`
* `seq.Map[string,int]{"one": 1, "two": 2}`
* `nil` is a valid empty map
* `make(seq.Map[string,int], 10)`
* `seq.Reduce()` with the `seq.MakeMap`, `seq.GroupBy`, or `seq.UpdateMap` collection functions

Map Features
----
Apart from being usable as a standard Go map, and a seq of tuples, maps have a few extra methods:

* `Map.Keys()` returns a `seq.Seq[K]`
* `Map.Values()` returns a `seq.Seq[V]`
* `Map.Contains(k)` returns a bool
* `Map.Get(k)` returns an `Opt[V]`