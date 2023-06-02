Slices
====
The `seq.Slice[T]` type is just a wrapper for `[]T` and all operations that can be done
on standard Go slices can be done on a `seq.Slice`. Eg

* `make(seq.Slice[string], 0, 10)`
* Literal assignment `mySlice := seq.Slice[string]{"one", "two"}`
* `len(mySlice)`
* for-range loops
* Indexing, like `one := mySlice[0]`
* Slicing `subSlice := mySlice[0:1]`

Limitations and Caveats
----
At the time of release of Fn 1.0 the Go compiler's support for generics is very
limited. This means that it can not tell that a `seq.Slice[T]` implements the `seq.Seq[T]`
interface. Either construct your sets with `seq.SliceOf()` which returns a `seq.Seq[T]` or
call `mySlice.Seq()` to explicitly cast it.

Creating Slices
----
A very short summary of different ways to create a Slice:

* `seq.SliceAs()`
* `seq.SliceAsArgs(...)`
* `seq.Slice[string]{"one", "two"}`
* `nil` is a valid empty slice
* `make(seq.Slice[string], 0, 10)`
* `seq.Reduce()` with the `seq.MakeSlice` collection function

Slice Features
----
The Slice struct adds a bunch of helpful methods on top of the standard Seq API. These include

 * Last()
 * One()
 * Copy()

And some ordering helpers that *manipulate the slice in-place*:
 * Sort()
 * Reverse()
 * Shuffle()

## The slice Package
Fn includes a package called 'slice' that works directly on standard Go
slices and maps and does not use seqs at all. The functions in 'slice' are intended
for doing one-shot conversions and mapping elements 1-1.

You can find a few [examples of how to use 'slice'](https://github.com/kamstrup/fn/blob/main/examples/slice_test.go)
in the "examples" folder.