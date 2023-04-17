Fn Documentation
====

This is the official documentation for the Fn package.

Overview
----
 * [Introduction](intro.md)
 * Core concepts:
   * [Seqs](seq.md)
   * [Opts](opt.md)
 * Working with:
   * [Slices](slice.md)
   * [Maps](map.md)
   * [Sets](set.md)
 * [Error Handling](error.md)
 * [Channels and Goroutines](parallel.md)
 * [Tips & Tricks](tips.md)
 * [API Reference](https://pkg.go.dev/github.com/kamstrup/fn) (external)

## Experimental Packages
The `seq`, `opt`, `slice`, `constraints`, and `fnmath` packages all have a stable API, but there are a few extra packages that
are not completely set in stone yet:

The API of the following packages are subject to change:

* `seqio` provides Seq[[]byte] based on io.Reader, and a Scanner based on bufio.Scanner
* `seqjson` provides Seq[T] based on json.Decoder
* `fntesting` contains various utilities to test Seqs