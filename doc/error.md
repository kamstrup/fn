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