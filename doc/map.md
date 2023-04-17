Working with Maps
====
One of the most important things to realize about maps in Fn is that `seq.Map[K,V]`
is just a type wrapper for `map[K]V`. This means that standard Go builtins and syntax
work on it:

* `make(seq.Map[string,int], 10)`
* Literal assignment `myMap := seq.Map[string,int]{"one": 1, "two": 2}`
* `len(myMap)`
* for-range loops
* Indexing, like `val, hasElement := myMap[key]`

Limitations and Caveats
----
At the time of release of Fn 1.0 the Go compiler's support for generics is very
limited. This means that it can not tell that a `seq.Map[K,V]` implements the `seq.Seq[Tuple[K,V]]`
interface. Either construct your maps with `seq.MapOf()` which returns a `seq.Seq[Tuple[K,V]]` or
call `myMap.Seq()` to explicitly cast it.

Creating Maps
----
A very short summary of different ways to create a map:

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

Examples
----

#### Building a Map with MakeMap and TupleWithKey
In order to use `MakeMap` to build a map we need a Seq of
`seq.Tuple`. The 2 easiest ways to obtain a Seq of tuples are
via mapping your seq with `TupleWithKey`, or via `ZipOf`.

This example uses `TupleWithKey` on a `*User` to build a
Seq of `Tuple[UserID, *User]` and collect that into a `map[UserID]*User`:
```go
type UserID uint64
type User struct { ID UserID ... }
usersSlice := []*User { ... }

users := seq.SliceOf(usersSlice)
userTuples := seq.MappingOf(users, seq.TupleWithKey(func(u *User) UserID {
   return u.ID
})
usersByID := seq.Reduce(nil, seq.MakeMap, userTuples).Or(nil)
// usersByID is a map[UserID]*User, the '.Or(nil)' above converts the Opt result to nil if there are errors
```

#### Counting Unique Names with UpdateMap
`UpdateMap` and `UpdateSlice` can be used in conjunction with
an *updater* function to create a collector. The updater function
tells the collector what to do if there is an existing value in a slot.

In this example we count the number of occurrences of names in a seq.
We do this by mapping to names onto a seq of `{name, 1}` tuples and then
instructing the `UpdateMap` to sum the values every time it merges an
element into the map:
```go
names := seq.SliceOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
tups := seq.ZipOf[string, int](names, seq.Constant(1))
res := seq.Reduce(nil, seq.UpdateMap[string, int](seq.Sum[int]), tups)
// res is an Opt[map[string,int]] with the value:
// map[string]int{
//   "bob":    3,
//   "alan":   2,
//   "scotty": 1,
// }
```
