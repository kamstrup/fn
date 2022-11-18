// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fn defines a set of useful constraints to be used
// with type parameters.
// This package was called "constriants" and was removed from the Go standard library
// before the 1.18 release, but is reproduced here inside the fn library because we
// need the numeric types. See https://github.com/golang/go/issues/50792
package fn

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
type Ordered interface {
	Integer | Float | ~string
}

// Arithmetic is the type union of all numeric types.
type Arithmetic interface {
	Integer | Float | Complex
}
