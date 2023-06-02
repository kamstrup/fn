Error Handling
====
When operating on in-memory structures like slices, maps, channels and so forth error handling is
normally not relevant. But if you do IO or some other operation that can error on runtime Fn provides
a few ways to handle it.

Operations that can produce errors are wrapped in `Opt[T]`.
Any error encountered via `sq.ForEach()`, `sq.ForEachIndex()`, `sq.First()` or `seq.Reduce()` are reported via opts.

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

### Aborting a Loop On First Error
You will normally abort a loop as soon as it encounters an internal error.
If you do not care for returning a result from the loop you can simply map
the inner operation directly to a Go `error` and check the result with `myLoop.First()`, like so
```go

var filesToDownload = seq.Slice[string]{"file1.txt", "file2.txt", "file3.txt"}

func downloadFile(name string) error {
	// store downloaded file directly in current working dir
}

func downloadAllFiles() error {
	firstErr, _ := seq.MappingOf(filesToDownload.Seq(), downloadFile).
		Where(seq.IsNonZero[error]).
		First()
	return firstErr.Or(nil)	
}
```