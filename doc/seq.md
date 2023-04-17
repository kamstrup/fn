Seqs
====
Seq is a core data structure of the Fn library. Seq is short for Sequence.

The Seq API is designed to work on top of immutable structures, where elements can be traversed sequentially.
An important thing to understand is that there is no stateful "iterator" on seqs.
Walking through a seq is done similarly to how you `append()` elements to a slice in Go, but inversely.

```go
ints = append(ints, i)
// OR if we call the slice "tail" and the element "head":
tail = append(tail, head)
```
The "inverse" of this operation looks like:
```go
head, tail = tail.First() // returns the first element of tail and a new tail
```
There are many easier ways to walk a Seq though. For example via `sq.ForEach()`
, `sq.Take()`, and `seq.Reduce()`. See [Iterating Over a Seq](#iterating-over-a-seq).

**Generally seqs are immutable**. Any exception to this will be clearly documented.

[Check the interface definition for Seq here](https://github.com/kamstrup/fn/blob/main/seq/seq.go).
There many ways to create seqs from standard Go structures.

### Creating Seqs
We follow the convention that functions for creating a Seq are named with an "Of"-suffix.
Ie `StringOf`, `SliceOf` etc. They always return a `Seq[T]`. Functions with an "As"-suffix
return a specific Seq implementation that allows you to perform type specific operations.
Fx. like sorting a `Slice`.

It is a known limitation of the Go compiler (v1.19) that it
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
Operations that iterate over a seq are said to *execute the seq*. This wording is intended to
signify that the operation might change the state of the seq. Although seqs can be assumed to be
immutable in many cases, that is not always the case. Examples of seqs that change state on execution are
channels, open files, or database connections.

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
head  := sq.Limit(10) // lazy seq of first 10 elements
mappedSameType := sq.Map(func (val T) T { ... }) // the Seq method Map() can only produce a seq of the same type
mappedOtherType := seq.MappingOf(seq, func(val T) S {}) // becomes a Seq[S]
```

### Transforming Seqs
Limiting the elements seen in a Seq is done with:
```go
sq.Where(predicate)
sq.While(predicate)
sq.TakeWhile(predicate) // if you also need the tail
sq.Limit(N)
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
The simplest way to collect results from a Seq is to call `sq.ToSlice()`.
It is often desirable to collect the elements into another structure that is not
just a slice. Maybe some sort of map, buffer, or completely custom data type.

To this end Fn has the functions `Reduce()`. This is also known as "fold" in
functional programming terminology.

#### Building a string with Reduce()
```go
strs := seq.SliceOfArgs("one", "two")
res := seq.Reduce(nil, seq.MakeString, strs)
// res is an Opt[strings.Builder] with the value "onetwo"
```

#### Collector Functions For seq.Reduce()
The first argument to `Reduce()` is a *collector function*.
Fn ships with a suite of standard collectors in the `seq` package, including:
`MakeSlice`, `MakeMap`, `MakeSet`, `MakeString`, `MakeBytes`, `Count`, and `GroupBy`.
For numeric properties we have
`fnmath.Sum`, `fnmath.Min`, `fnmath.Max`, and `fnmath.MakeStats`.

There are 2 more advanced collection helpers `UpdateMap`, `UpdateSlice`.


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