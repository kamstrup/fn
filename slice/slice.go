package slice

import (
	"sort"

	"github.com/kamstrup/fn/constraints"
)

// Mapping converts a slice from one type to another
func Mapping[S, T any](slice []S, f func(s S) T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := make([]T, len(slice))
	for i, s := range slice {
		result[i] = f(s)
	}
	return result
}

// MappingIndex converts a slice from one type to another.
// The mapping function also receives the index of the element being transformed.
func MappingIndex[S, T any](slice []S, f func(i int, s S) T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := make([]T, len(slice))
	for i, s := range slice {
		result[i] = f(i, s)
	}
	return result
}

// Gen builds a new slice of a given size.
// The generator function is called for each index in the new slice.
func Gen[T any](sz int, generator func(i int) T) []T {
	result := make([]T, sz)
	for i := range result {
		result[i] = generator(i)
	}
	return result
}

// Copy returns a shallow copy of the given slice.
func Copy[T any](slice []T) []T {
	cpy := make([]T, len(slice))
	copy(cpy, slice)
	return cpy
}

// Zero zeroes all elements in a slice.
func Zero[T any](slice []T) []T {
	// compiles to runtime.memclr(), see https://github.com/golang/go/issues/5373
	var zero T
	for i := range slice {
		slice[i] = zero
	}
	return slice
}

// SortAsc sorts a slice in-place.
// The argument is returned to facilitate easy chaining.
func SortAsc[T constraints.Ordered](slice []T) []T {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	return slice
}

// SortDesc sorts a slice in-place.
// The argument is returned to facilitate easy chaining.
func SortDesc[T constraints.Ordered](slice []T) []T {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] > slice[j]
	})
	return slice
}
