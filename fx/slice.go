package fx

import (
	"sort"

	"github.com/kamstrup/fn/constraints"
)

// MapSlice converts a slice from one type to another
func MapSlice[S, T any](slice []S, f func(s S) T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := make([]T, len(slice))
	for i, s := range slice {
		result[i] = f(s)
	}
	return result
}

// MapSliceIndex converts a slice from one type to another.
// The mapping function also receives the index of the element being transformed.
func MapSliceIndex[S, T any](slice []S, f func(i int, s S) T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := make([]T, len(slice))
	for i, s := range slice {
		result[i] = f(i, s)
	}
	return result
}

// GenSlice builds a new slice of a given size.
// The generator function is called for each index in the new slice.
func GenSlice[T any](sz int, generator func(i int) T) []T {
	result := make([]T, sz)
	for i := range result {
		result[i] = generator(i)
	}
	return result
}

// CopySlice returns a shallow copy of the given slice.
func CopySlice[T any](slice []T) []T {
	cpy := make([]T, len(slice))
	copy(cpy, slice)
	return cpy
}

// ZeroSlice zeroes all elements in a slice.
func ZeroSlice[T any](slice []T) []T {
	// compiles to runtime.memclr(), see https://github.com/golang/go/issues/5373
	var zero T
	for i := range slice {
		slice[i] = zero
	}
	return slice
}

// SortSliceAsc sorts a slice in-place.
// The argument is returned to facilitate easy chaining.
func SortSliceAsc[T constraints.Ordered](slice []T) []T {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	return slice
}

// SortSliceDesc sorts a slice in-place.
// The argument is returned to facilitate easy chaining.
func SortSliceDesc[T constraints.Ordered](slice []T) []T {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] > slice[j]
	})
	return slice
}

// AssocSlice builds a map[K]V from a slice.
// To build a map from a slice use SliceAssoc.
// The mnemonic for the function name is that the last word is what you operate on.
func AssocSlice[K comparable, V, T any](slice []T, f func(T) (K, V)) map[K]V {
	m := make(map[K]V, len(slice))
	for _, t := range slice {
		k, v := f(t)
		m[k] = v
	}
	return m
}
