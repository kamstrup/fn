Error Handling
====
When operating on in-memory structures like slices, maps, channels and so forth error handling is
normally not relevant. But if you do IO or some other operation that can error on runtime Fn provides
a few ways to handle it.

The `seq.Error(seq)` function returns an `error` if there is an error associated with a Seq or Opt.
When you execute a Seq the "empty" tail Seq from `ForEach()` and other operations will capture any
errors.

Alternatively you can wrap results in `Opt[T]` which can also capture an error.
Any error encountered via `sq.First()` or `seq.Reduce()` are reported via opts.

#### Warning
When calling `Seq.ToSlice()` all error are ignored. ToSlice is intended for in-memory operations.

Working with Functions that Return Errors
----
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
optInt1 := opt.Returning(getInt())
optInt2 := opt.Recovering(calcInt(27)) // to recover panics
```