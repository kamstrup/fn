// Package slice is a lightweight package of functions for transforming standard
// in-memory Go structures.
//
// This package leans into the terminology and patterns from "seq" package,
// but differs in a few key areas:
//
//   - All functions are eager, i.e. executed where they are used. Nothing is lazy.
//   - Designed to do one-time, and one-to-one, transformations
//   - There is no error handling. It is designed to work with in-memory transformations,
//     not IO or other external resources.
//
// If you need to do more complex data transformation, like filtering out subsets of elements,
// handling errors, or any longer chain of transformations, you will be better off by using
// seq.Seq.
//
// # Interoperability with seq.Slice
//
// It is important to note that all functions that take a []T argument can be passed a seq.Slice[T]
// as well. Any seq.Slice from the seq package can be passed as arguments to functions in this package
// (and anywhere else a slice is used).
package slice
