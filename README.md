# Fn(), Functional Programming for Golang

[![Go Report Card](https://goreportcard.com/badge/github.com/kamstrup/fn)](https://goreportcard.com/report/github.com/kamstrup/fn) [![PkgGoDev](https://pkg.go.dev/badge/github.com/kamstrup/fn)](https://pkg.go.dev/github.com/kamstrup/fn)

Fn is library for golang that enable you to blend functional programming techniques
with standard idiomatic Go code.

We are inspired by Clojure and the Java Streams APIs that were
introduced back in Java 8, and want to provide something of similar spirit
that makes it even more fun to write Go code.

## Philosophy
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

## Terminology
For starters let's get some terminology in place.

**Seq:** The core data structure is `Seq[T]`. It is short for Sequence.
The Seq API is designed to work on top of immutable structures,
thus there is no stateful "iterator". Walking through a Seq is done similarly 
to how you `append()` elements to a slice in Go, but inversely.

```go
ints = append(ints, i)
// OR if we call the slice "tail" and the element "head":
tail = append(tail, head)
```
The "inverse" of this operation looks like:
```go
head, tail = tail.First() // pops the first element and returns a new tail
```
There are many easier ways to walk a Seq though. For example via `seq.ForEach()`
, `seq.Take()`, and `fn.Into()`. See [Iterating Over a Seq](#iterating-over-a-seq).

Generally seqs are immutable. Any exception to this will be clearly documented. 

[Check the interface definition for Seq here](https://github.com/kamstrup/fn/blob/main/seq.go).
There many ways to create seqs from standard Go structures. You can find most of them in 
the section [Creating Seqs](#creating-seqs).

**Array:** Standard Go slices `[]T` are wrapped in the `fn.Array[T]` type.
The `Array` type is a public subtype of `[]T` so you can do numeric indexing on an `Array`.
Arrays are seqs, but also add some extra methods like `Sort()` and `Reverse()`.

**Assoc:** Because "map" is an overloaded word in functional programming,
Fn() uses the word "assoc" instead of `map[K]V` (pronounced with a soft "ch" at the end).
The word "map" is reserved for the mapping operation used to convert a Seq to something else.

**Tuple:** Or "pair". Represents two data points. A helper mainly used when working with assocs,
where the tuple captures a key and a value. Assocs can be interpreted as a seq of tuples, or if you
build a seq of tuples you can create an assoc or map from it.

**Opt:** Returned from operations where you are not certain to get a result.
For example when you call `seq.First()`. If the seq is empty you get back an empty opt,
and an empty tail seq.

**Seq.Len():** Lengths are handled in a special way in Fn(). They are allowed to be finite, 
unknown, or infinite. Making these distinctions opens the possibility of pre-allocating
slices and maps of the correct size, which can make a big difference in performance critical code.

## API Overview
If you just want to jump in and see some code you can check out
[the simple examples](https://github.com/kamstrup/fn/blob/main/examples/simple_test.go).
Otherwise here follows a brief overview.

Fn also bundles a [very simple sub-library called Fx](#the-fx-package---simplified-fn)
that you can use to do 1-line functional constructs.

### Creating Seqs
We follow the convention that functions for creating a Seq are named with an "Of"-suffix.
Ie `StringOf()`, `ArrayOf` etc. They always return a `Seq[T]`. Functions with an "As"-suffix
return a specific Seq implementation that allows you to perform type specific operations.
Fx. like sorting an `Array`. It is a known limitattion of the Go compiler (v1.19) that it 
can not determine that generic structs implement generic interfaces. So in order to use
an `Array`, `Assoc`, `Set`, or `String` as a Seq you need to call `.Seq()` on the instance.
Seq creation funcs that take a variadic list of arguments have an "OfArgs"-suffix.

#### From Standard Go Types
```go
arr := fn.ArrayOfArgs(1,2,3) // also: ArrayOf(), ArrayAs(), ArrayAsArgs()
ass := fn.AssocOf(map[string]int{"foo": 27, "bar": 68}) // also: AssocAs()
set := fn.SetOf(map[string]struct{}{"foo", {}, "bar": {}}) // also: SetAs()
str := fn.StringOf("hello world") // also: StringAs()
ch := fn.ChanOf(make(chan T))
twentySeven := fn.SingletOf(27) // single element Seq
empty := fn.EmptySeq[int]()
```

#### Numeric Ranges
```go
zero := fn.Constant(0) // infinite
nums := fn.RangeOf(0, 10)
evenNums := fn.RangeStepOf(0, 10, 2)
toInfinity := fn.NumbersFrom(0) 
```

#### From Functions or Closures
```go
src := fn.SourceOf(func T { ... }) // infinite
```

### Iterating over a Seq
Functions that execute the Seq, ie actively traverse it include:
```go
seq.ForEach(func(elem T) {
   // use elem
})
seq.ForEachIndex(func (i int, elem T) {
   // use index and elem
})
tenFirst, tailSeq := seq.Take(10)
goodArray, tailSeq := seq.TakeWhile(func(elem T) bool { return isGood(elem)})
```
There are also some helper functions included in Fn for executing a Seq for various purposes.
See the [Operations on Seqs](#operations-on-seqs).

Functions that do not execute the Seq, but return a new lazy Seq include:
```go
while := seq.While(predicate)
where := seq.Where(predicate)
mappedSameType := seq.Map(func (val T) T { ... }) // the Seq method Map() can only produce a seq of the same type
mappedOtherType := fn.MapOf(seq, func(val T) S {}) // becomes a Seq[S]
```

### Transforming Seqs
Limiting the elements seen in a Seq is done with:
```go
seq.Where(predicate)
seq.While(predicate)
seq.TakeWhile(predicate) // if you also need the tail
```
Transforming elements, mapping them 1-1 is done with
```go
seq.Map(func(t T) T { ... })
seqT := fn.MapOf(seqS, func(s S) T { ... })
```
You can split a `Seq[T]` into sub-seqs with
```go
subs := fn.SplitOf(seq, splitterFunc)
```
and you can join seqs together with
```go
longSeq := fn.ConcatOf(seq1, seq2, ... )
longSeq := fn.FlattenOf(seqOfSeqs)
```

If you have 2 seqs that you want to traverse in parallel as pairs of elements
you can use `ZipOf`:
```go
ints := fn.ArrayOfArgs(1,2,3)
strs := fn.ArrayOfArgs("one", "two", "three")
pairs := fn.ZipOf(ints, strs)
// pairs is a Seq[Tuple[int,string]]
```

### Predicates
Predicates that can be used directly on any ordered type T:
```go
fn.IsZero[T] // matching the zero value of a type T 
fn.IsNonZero[T] // matching any non-zero value of a type T
fn.GreaterThanZero[T] // > zero values for T
fn.LessThanZero[T] // < zero value for T
```

Functions that can help you create a predicate:
```go
fn.Is(x) // val == x
fn.IsNot(x) // val != x
fn.Not(pred) // !pred(val)
fn.GreaterThan(x) // val > x
fn.LessThan(x) // val < x
```

### Collecting Results
The simplest way to collect results from a Seq is to call `seq.Array()`.
It is often desirable to collect the elements into another structure that is not
just a slice. Maybe some sort of map, buffer, or completely custom data type.

To this end Fn has the functions `Into()` and `IntoErr()`. This is often
also known as "reduce" or "fold" in functional programming terminology.

#### Building a string with Into()
```go
strs := fn.ArrayOfArgs("one", "two")
res := fn.Into(nil, fn.MakeString, strs)
// res is a Opt[string] with the value "onetwo"
```

#### Collector Functions For fn.Into()
The second argument to `Into()` is a *collector function*.
Fn ships with a suite of standard collectors including:
`Append`, `MakeAssoc`, `MakeSet`, `MakeString`, `MakeBytes`,
`Sum`, `Count`, `Min`, `Max`, and `GroupBy`. There
are 2 more advanced collection helpers `UpdateAssoc`, `UpdateArray`.

#### Building a Map with MakeAssoc and TupleWithKey
In order to use `MakeAssoc` to build a map we need a Seq of
`fn.Tuple`. The 2 easiest ways to obtain a Seq of tuples are
via mapping your seq with `TupleWithKey`, or via `ZipOf`.

This example uses `TupleWithKey` on a `*User` to build a
Seq of `Tuple[UserID, *User]` and collect that into a `map[UserID]*User`:
```go
type UserID uint64
type User struct { ID UserID ... }
usersSlice := []*User { ... }

users := fn.ArrayOf(usersSlice)
userTuples := fn.MapOf(users, fn.TupleWithKey(u *User) UserID {
   return u.ID
})
usersByID := fn.Into(nil, fn.MakeAssoc, userTuples).Or(nil)
// usersByID is a map[UserID]*User, the '.Or(nil)' above converts the Opt result to nil if there are errors
```

#### Counting Unique Names with UpdateAssoc
`UpdateAssoc` and `UpdateArray` can be used in conjunction with
an *updater* function to create a collector. The updater function
tells the collector what to do if there is an existing value in a slot.

In this example we count the number of occurrences of names in a Seq.
We do this by mapping to names onto a seq of `{name, 1}` tuples and then
instructing the `UpdateAssoc` to sum the values every time it merges an
element into the map:
```go
names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
tups := fn.ZipOf[string, int](names, fn.Constant(1))
res := fn.Into(nil, fn.UpdateAssoc[string, int](fn.Sum[int]), tups)
// res is an Opt[map[string,int]] with the value:
// map[string]int{
//   "bob":    3,
//   "alan":   2,
//   "scotty": 1,
// }
```

### Operations on Seqs
To check if a Seq contains some given element you can use `fn.Any(seq, pred)`:
```go
nums := fn.RangeOf(0, 10)
hasEvenNum := fn.Any(nums, func (n int) bool { return n % 2 == 0})
hasSeven := fn.Any(nums, fn.Is(7))
```
You can also check if all elements satisfy some criteria with `fn.All(seq, pred)`.

Similar to how you can retrieve the first element in a Seq with `head, tail := seq.First()`
you can get the last element with `last := fn.Last(seq)`.

Executing a Seq for side effects, fx. printing all elements, can be done with `fn.Do()`:
```.go
nums := fn.RangeOf(0, 10).
   Map(func (n int) int {
      fmt.Println(n)
      return n
   })

// Nothing is printed since 'nums' is lazy.
// We can force it to execute with:
fn.Do(nums)
// prints numbers from [0..9]
```

### Opts
Some operations return `Opt[T]`, notably `seq.First()` and `Into()`.
Opts are used to represent a value that might not be there (if the seq is empty), or capture potential errors.
An opt with a captured error is considered empty, and empty opts will report the error `fn.ErrEmpty`.

They have a range of helper API that allows for easy chaining:
```go
opt.Must()   // Returns T or panics if the if the opt is empty
opt.Empty()  // True if there was an error or no value is captured
opt.Ok()     // Opposite of Empty()
opt.Return() // Unpacks into standard "T, error" pair
opt.Or(val)  // Returns val if opt is empty, or the option's own value if non-empty

opt.Map(func(val T) T { ... }) // Converts the value to something of the same type
OptMap(opt, func(val T) S)     // Converts the opt into another type

opt.Seq() // Interprets the option as a single-valued Seq

opt.Error() // Returns nil, fn.ErrEmpty, or any captured error
opt.OnErr(func (error) T { ... }) // Returns the opt value T, or invokes a callback with the error
```

**Opt Misconceptions and Pitfalls**: Opts should be passed by value, on the stack.
If you use pointers, `*Opt` or `&opt`, something is wrong. They have a bit of memory overhead
and should not be stored in arrays or slices. You *can* use them in seqs because they are
lazily created 1 by 1 and kept short-lived on the stack.
An Opt is not a "promise" or "future" - they capture an existing result. 


### Parallel Execution
You can execute a Seq in N goroutines mapping the results into a new Seq with `fn.Go()`:
```.go
func fetchItem(id int) Opt[T] {
   // do something slow and calculate t
   return fn.OptOf(t) // or maybe an error
}

// Execute fetchItem of 1027 ids in 100 parallel goroutines
ids := fn.RangeOf(0, 1027)
result := fn.Go(ids, 100, fetchItem)

// result is a Seq[Opt[T]], let's print the successes and errors 
result.ForEach(func (opt Opt[T]) {
   t, err := opt.Error()
   if err != nil {
      fmt.Println("Oh no, an error!", err)
   } else {
      fmt.Println("Nice, got one T:", t)
   }
})
```

### Error Handling
When operating on in-memory structures like slices, maps, channels and so forth error handling is
normally not relevant. But if you do IO or some other operation that can error on runtime Fn provides
a few ways to handle it.

The `fn.Error(seq)` function returns an `error` if there is an error associated with a Seq or Opt.
When you execute a Seq the "empty" tail Seq from ForEach() and other operations will capture any
errors.

Alternatively you can wrap results in `Opt[T]` which can also capture an error.
Any error encountered via `seq.First()` or `fn.Into()` are reported via opts.

## The Fx Package - Simplified Fn
Fn includes a minimal sub-library called Fx that works directly on standard Go
slices and maps and does not use seqs at all. The functions in Fx are intended
for doing one-shot conversions and mapping elements 1-1.

You can find a few [examples of how to use Fx](https://github.com/kamstrup/fn/blob/main/examples/fx_test.go)
in the "examples" folder.
.

## Performance
If the foundational functional data structures and algorithms is not done carefully,
execution speed and memory usage will suffer. Fn() is designed to make the best of what
the Go runtime can provide. Initial benchmarks puts it as a top contender among Golang
functional libraries. See benchmarks here https://github.com/mariomac/go-stream-benchmarks/pull/1

TODO
---
```
API CHANGES:
* Do we need to change fn.Go() to enable better error handling?

POTENTIAL FUTURE FEATURES (in order of prio)
* fnio.DirOf(dirName), * fnio.DirTreeOf(dirName) (recursive)
* Special seqs for Assoc.Keys() and Assoc.Values()
* seq.Limit(n) Seq[T], lazy counterpart to seq.Take(n)
* RunesOf(string) Seq[rune]
* MakeChan collector func for Into()?
* A small JSON package "fnjson" to help reading and writing Seqs of JSON objects
* MultiChan() Seq that selects on multiple chan T?
* Something for context.Context? Support cancel() cb and Done() chans? fncontext package...
* fn.GoErr(seq, numTasks, FuncMapErr) -- or some version of fn.Go() with cancellation and error handling. 
* Tuple[S,T] as Seq[any]? (we have to do "any" bc the types S!=T)
* MergeSort[T any](FuncLess[T], seqs ... Seq[T]) Seq[T] -- lazy merge sorting of pre-sorted Seqs
* Compound FuncCollect, CollectorOf[S,T any](funcs ... FuncCollect[S,T]) FuncCollect[S,[]T]
* Seq[Arithmetic] producing random numbers?
* Seq for *sql.Rows, with some type safe mechanism for reading rows
* Promises or Futures that work nicely with Seq and Opt?

POTENTIAL FUTURE OPTIMIZATIONS
* Seq of single element (see SingletOf(t))
* EmptySeq impl. (currently just wraps an empty slice), but an empty struct{} would do even better
* Look for allocating buffers of right size where we can
* Can we do some clever allocations in fn.Into() when seed is nil?
```