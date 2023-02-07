# Fn(), Functional Programming for Golang

[![Go Report Card](https://goreportcard.com/badge/github.com/kamstrup/fn)](https://goreportcard.com/report/github.com/kamstrup/fn) [![PkgGoDev](https://pkg.go.dev/badge/github.com/kamstrup/fn)](https://pkg.go.dev/github.com/kamstrup/fn)

Fn is library for golang that enable you to blend functional programming techniques
with standard idiomatic Go code.

We are inspired by [Clojure](https://github.com/clojure/clojure),
[Vavr](https://github.com/vavr-io/vavr), and the Java Streams APIs that were
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
There are many easier ways to walk a Seq though. For example via `sq.ForEach()`
, `sq.Take()`, and `seq.Reduce()`. See [Iterating Over a Seq](#iterating-over-a-seq).

Generally seqs are immutable. Any exception to this will be clearly documented. 

[Check the interface definition for Seq here](https://github.com/kamstrup/fn/blob/main/seq/seq.go).
There many ways to create seqs from standard Go structures. You can find most of them in 
the section [Creating Seqs](#creating-seqs).

**Slice:** Standard Go slices `[]T` are wrapped in the `seq.Slice[T]` type.
The `Slice` type is a public subtype of `[]T` so you can do numeric indexing on a `Slice`,
and use `cap()`, `len()`, and for-range loops on them.
Slices are seqs, but also add some extra methods like `Sort()` and `Reverse()`.

**Tuple:** Or "pair". Represents two data points. A helper mainly used when working with maps,
where the tuple captures a key and a value. Maps can be interpreted as a seq of tuples, or if you
build a seq of tuples you can create map from it.

**Opt:** Returned from operations where you are not certain to get a result.
For example when you call `sq.First()`. If the seq is empty you get back an empty opt,
and an empty tail sq.

**sq.Len():** Lengths are handled in a special way in Fn(). They are allowed to be finite, 
unknown, or infinite. Making these distinctions opens the possibility of pre-allocating
slices and maps of the correct size, which can make a big difference in performance critical code.
In most circumstances you will not need to use Len(). All operations are valid on empty seqs and
empty opts, so just set up your pipeline of operations and check if the end result is valid.

## API Overview
If you just want to jump in and see some code you can check out
[the simple examples](https://github.com/kamstrup/fn/blob/main/examples/simple_test.go).
Otherwise here follows a brief overview.

Fn also bundles a [very simple sub-library called 'slice'](#the-slice-package)
that you can use to do 1-line functional constructs.

### Creating Seqs
We follow the convention that functions for creating a Seq are named with an "Of"-suffix.
Ie `StringOf()`, `SliceOf` etc. They always return a `Seq[T]`. Functions with an "As"-suffix
return a specific Seq implementation that allows you to perform type specific operations.
Fx. like sorting a `Slice`. It is a known limitation of the Go compiler (v1.19) that it 
can not determine that generic structs implement generic interfaces. So in order to use
a `Slice`, `Map`, `Set`, or `String` as a Seq you need to call `.Seq()` on the instance.
Seq creation funcs that take a variadic list of arguments have an "OfArgs"-suffix.

#### From Standard Go Types
```go
arr := seq.SliceOfArgs(1,2,3) // also: SliceOf(), SliceAs(), SliceAsArgs()
ass := seq.MapOf(map[string]int{"foo": 27, "bar": 68}) // also: MapAs()
set := seq.SetOf(map[string]struct{}{"foo", {}, "bar": {}}) // also: SetAs()
str := seq.StringOf("hello world") // also: StringAs()
ch := seq.ChanOf(make(chan T))
twentySeven := seq.SingletOf(27) // single element Seq
empty := seq.Empty[int]()
```

#### Numeric Ranges
```go
zero := seq.Constant(0) // infinite
nums := seq.RangeOf(0, 10)
evenNums := seq.RangeStepOf(0, 10, 2)
toInfinity := seq.RangeFrom(0) // "infinity" == max value for the numeric type 
```

#### From Functions or Closures
```go
src := seq.SourceOf(func T { ... }) // infinite
```

### Iterating over a Seq
Functions that execute the Seq, ie actively traverse it include:
```go
sq.ForEach(func(elem T) {
   // use elem
})
sq.ForEachIndex(func (i int, elem T) {
   // use index and elem
})
tenFirst, tailSeq := sq.Take(10)
goodArray, tailSeq := sq.TakeWhile(func(elem T) bool { return isGood(elem)})
```
There are also some helper functions included in Fn for executing a Seq for various purposes.
See the [Operations on Seqs](#operations-on-seqs).

Functions that do not execute the Seq, but return a new lazy Seq include:
```go
while := sq.While(predicate)
where := sq.Where(predicate)
mappedSameType := sq.Map(func (val T) T { ... }) // the Seq method Map() can only produce a seq of the same type
mappedOtherType := seq.MappingOf(seq, func(val T) S {}) // becomes a Seq[S]
```

### Transforming Seqs
Limiting the elements seen in a Seq is done with:
```go
sq.Where(predicate)
sq.While(predicate)
sq.TakeWhile(predicate) // if you also need the tail
```
Transforming elements, mapping them 1-1 is done with
```go
sq.Map(func(t T) T { ... })
seqT := seq.MappingOf(seqS, func(s S) T { ... })
```
You can split a `Seq[T]` into sub-seqs with
```go
subs := seq.SplitOf(seq, splitterFunc)
```
and you can join seqs together with
```go
longSeq := seq.ConcatOf(seq1, seq2, ... )
longSeq := seq.FlattenOf(seqOfSeqs)
longSeq := seq.Prepend(value, seq1) // prepends a single value to a Seq
```

If you have 2 seqs that you want to traverse in parallel as tuples (pairs) of elements
you can use `ZipOf`:
```go
ints := seq.SliceOfArgs(1,2,3)
strs := seq.SliceOfArgs("one", "two", "three")
pairs := seq.ZipOf(ints, strs)
// pairs is a Seq[Tuple[int,string]]
```

### Predicates
Predicates that can be used directly on any ordered type T:
```go
seq.IsZero[T] // matching the zero value of a type T 
seq.IsNonZero[T] // matching any non-zero value of a type T
seq.GreaterThanZero[T] // > zero values for T
seq.LessThanZero[T] // < zero value for T
```

Functions that can help you create a predicate:
```go
seq.Is(x) // val == x
seq.IsNot(x) // val != x
seq.Not(pred) // !pred(val)
seq.GreaterThan(x) // val > x
seq.LessThan(x) // val < x
```

### Collecting Results
The simplest way to collect results from a Seq is to call `sq.Values()`.
It is often desirable to collect the elements into another structure that is not
just a slice. Maybe some sort of map, buffer, or completely custom data type.

To this end Fn has the functions `Reduce()`. This is also known as "fold" in
functional programming terminology.

#### Building a string with Reduce()
```go
strs := seq.SliceOfArgs("one", "two")
res := seq.Reduce(nil, seq.MakeString, strs)
// res is an Opt[string] with the value "onetwo"
```

#### Collector Functions For seq.Reduce()
The first argument to `Reduce()` is a *collector function*.
Fn ships with a suite of standard collectors in the `seq` package, including:
`Append`, `MakeMap`, `MakeSet`, `MakeString`, `MakeBytes`,
`fnmath.Sum`, `Count`, `fnmath.Min`, `fnmath.Max`, and `GroupBy`. There
are 2 more advanced collection helpers `UpdateMap`, `UpdateSlice`.

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
userTuples := seq.MappingOf(users, seq.TupleWithKey(u *User) UserID {
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

### Operations on Seqs
To check if a Seq contains some given element you can use `seq.Any(sq, pred)`:
```go
nums := seq.RangeOf(0, 10)
hasEvenNum := seq.Any(nums, func (n int) bool { return n % 2 == 0}) // true
hasSeven := seq.Any(nums, seq.Is(7)) // true
allNonZero := seq.All(nums, seq.IsNonZero[int]) // false
// Note: Even if all 3 calls to Any, Any, and All above execute the nums seq
//       this all still work because range-seqs are stateless.
```
You can also check if all elements satisfy some criteria with `seq.All(sq, pred)`.

Similar to how you can retrieve the first element in a Seq with `head, tail := sq.First()`
you can get the last element with `last := seq.Last(sq)`.

You can check if a seq is empty with `seq.IsEmpty(sq)`.

Executing a Seq for side effects, fx. printing all elements, can be done with `seq.Do()`:
```.go
nums := seq.RangeOf(0, 10).
   Map(func (n int) int {
      fmt.Println(n)
      return n
   })

// Nothing is printed since 'nums' is lazy.
// We can force it to execute with:
seq.Do(nums)
// prints numbers from [0..9]
```

### Opts
Some operations return `Opt[T]`, notably `sq.First()` and `seq.Reduce()`.
Opts are used to represent a value that might not be there (if the seq is empty), or capture potential errors.
An opt with a captured error is considered empty, and empty opts will report the error `opt.ErrEmpty`.

They have a range of helper API that allows for easy chaining:
```go
op.Must()   // Returns T or panics if the opt is empty
op.Empty()  // True if there was an error or no value is captured
op.Ok()     // Opposite of Empty()
op.Return() // Unpacks into standard "T, error" pair
op.Or(val)  // Returns val if opt is empty, or the option's own value if non-empty
op.Map(func(val T) T { ... }) // Converts the value to something of the same type
op.Error()  // Returns nil, opt.ErrEmpty, or any captured error
op.OnErr(func (error) T { ... }) // Returns the opt value T, or invokes a callback with the error

// And the static function:
opt.Map(opt, func(val S) T) Opt[T]  // Converts the opt from type S to T
```

**Opt Misconceptions and Pitfalls**: Opts should be passed by value, on the stack.
If you use pointers, `*Opt` or `&opt`, something is wrong. They have a bit of memory overhead
and should not be stored in arrays or slices. You *can* use them in seqs because they are
lazily created 1 by 1 and kept short-lived on the stack.
An Opt is not a "promise" or "future" - they capture an existing result. 


### Parallel Execution
You can execute a Seq in N goroutines mapping the results into a new Seq with `seq.Go()`:
```.go
func fetchItem(id int) Opt[T] {
   // do something slow and calculate t
   return seq.OptOf(t) // or maybe an error
}

// Execute fetchItem of 1027 ids in 100 parallel goroutines
ids := seq.RangeOf(0, 1027)
result := seq.Go(ids, 100, fetchItem)

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

The `seq.Error(seq)` function returns an `error` if there is an error associated with a Seq or Opt.
When you execute a Seq the "empty" tail Seq from ForEach() and other operations will capture any
errors.

Alternatively you can wrap results in `Opt[T]` which can also capture an error.
Any error encountered via `sq.First()` or `seq.Reduce()` are reported via opts.

## The slice Package
Fn includes a package called 'slice' that works directly on standard Go
slices and maps and does not use seqs at all. The functions in 'slice' are intended
for doing one-shot conversions and mapping elements 1-1.

You can find a few [examples of how to use 'slice'](https://github.com/kamstrup/fn/blob/main/examples/slice_test.go)
in the "examples" folder.

## Tips and Tricks

### Don't Check Length or Presence Until the Last step
All of the code in Fn() works with nil slices and maps, and empty Opts.

For example, to ensure that there is one and only one record with a given ID in the 10 first records
in some slice:
```.go
recordID := 1234
recSeq := seq.SliceOf(records).
    Take(10).
    Where(func(rec *Record) bool { return rec.ID == recordID})

theOneRecord, err := seq.One(recSeq).Return()
```
Note how we are not checking the length of 'records', or how many records with the given ID
we found. That is all handled by `seq.One()`.

### Chan, Map, Set, Slice and String Can be Used As Their Native Go Types
All the seq constructors names with the "As" suffix return their type wrapper.
```.go
myMap := seq.MapAs(make(map[string]int))

// We can do normal map[string]int stuff
myMap["one"] = 1
myMapLen := len(myMap)
for k, v := range myMap { ... }

// But also treat it as a seq
result := myMap.Where(...).ToSlice()
```

### Working with Functions that Return Errors
It is very common in Go to have function that look like
`func getInt() (int, error)` or `func calcInt(n int) (int, error)`.
This can make function chaining clumsy because you cannot pass
the results directly into another function.
In Fn this can helped out with opts. Functions that return errors
can be wrapped as functions returning opts:
```go
caller := opt.Caller(getInt) // is a func() Opt[int]
mapper := opt.Mapper(calcInt) // is a func(int) Opt[int]
```
These functions `caller`, and `mapper`, can be plugged directly into `seq.MappingOf()`, `seq.SourceOf()`,
and many others.

Or if you want to calculate the result immediately
```go
optInt1 := opt.Call(getInt)
optInt2 := opt.Apply(calcInt, 27)
```
There are panic-recovering variations of these functions as well.

## Performance
If the foundational functional data structures and algorithms is not done carefully,
execution speed and memory usage will suffer. Fn() is designed to make the best of what
the Go runtime can provide. Initial benchmarks puts it as a top contender among Golang
functional libraries. See benchmarks here https://github.com/mariomac/go-stream-benchmarks/pull/1

## Experimental Packages
 * `seqio` provides Seq[[]byte] based on io.Reader, and a Scanner based on bufio.Scanner
 * `seqjson` provides Seq[T] based on json.Decoder
 * `fntesting` contains various utilities to test Seqs

TODO
---
```
API CHANGES:
* Do we need to change seq.Go() to enable better error handling?
* Add Seq.Limit(n) method on Seq interface?

POTENTIAL FUTURE FEATURES (unordered)
* Something for context.Context? Support cancel() cb and Done() chans? fncontext package...
* seqio.DirOf(dirName), seqio.DirTreeOf(dirName) (recursive)
* Special seqs for Map.Keys() and Map.Values()
* RunesOf(string) Seq[rune]
* MakeChan collector func for Reduce()?
* MultiChan() Seq that selects on multiple chan T?
* MergeSort[T any](FuncLess[T], seqs ... Seq[T]) Seq[T] -- lazy merge sorting of pre-sorted Seqs
* Compound FuncCollect, CollectorOf[S,T any](funcs ... FuncCollect[S,T]) FuncCollect[S,[]T]
* Seq[Arithmetic] producing random numbers (in fnmath)?
* Seq for *sql.Rows, with some type safe mechanism for reading rows
* Promises or Futures that work nicely with Seq and Opt?
* Some kind of "push seq", or is that just Chan? Some libraries only provide "callback based iteration" for data structures.

POTENTIAL FUTURE OPTIMIZATIONS
* EmptySeq impl. (currently just wraps an empty slice), but an empty struct{} would do even better
* Look for allocating buffers of right size where we can
* Can we do some clever allocations in seq.Reduce() when seed is nil?
```