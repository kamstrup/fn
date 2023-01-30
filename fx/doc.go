// Package fx is a lightweight package of functions for transforming standard
// in-memory Go structures.
//
// This package leans into the terminology and patterns from the parent "fn" package,
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
package fx
