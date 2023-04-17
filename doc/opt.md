Opts
====

Some operations return `Opt[T]`, notably `sq.First()` and `seq.Reduce()`.
Opts are used to represent a value that might not be there (if the seq is empty), or capture potential errors.
An opt with a captured error is considered empty, and empty opts will report the error `opt.ErrEmpty`.

Opts have a range of helper API that allows for easy chaining:
```go
op.Empty()  // True if there was an error or no value is captured
op.Ok()     // Opposite of Empty()
op.Return() // Unpacks into standard "T, error" pair
op.Or(val)  // Returns val if opt is empty, or the option's own value if non-empty
op.Map(func(val T) T { ... }) // Converts the value to something of the same type
op.Error()  // Returns nil, opt.ErrEmpty, or any captured error
op.OnErr(func (error) T { ... }) // Returns the opt value T, or invokes a callback with the error
op.Must()   // Returns T or panics if the opt is empty

// And the static function:
opt.Map(opt, func(val S) T) Opt[T]  // Converts the opt from type S to T

// For async results
opt.Promise(...) Future[T]
```

**Opt Misconceptions and Pitfalls**: Opts should be passed by value, on the stack.
If you use pointers, `*Opt` or `&Opt`, something is wrong. They have a bit of memory overhead
and should not be stored in arrays or slices. You *can* use them in seqs because they are
lazily created 1 by 1 and kept short-lived on the stack.

An Opt is not a "promise" or "future" - they capture an existing result. If you need async
operations please look at `opt.Promise()`.

Seqs of Opts
----
When handling errors you often end up with a `Seq[Opt[T]]`. If you need to convert that into
a `Seq[T]` that can be done with `seq.ValuesOf()` helper.

