package slice

import (
	"math/rand"
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
// The generator function is called sz times.
//
// Example:
//
//	tenRandomFloats := slice.Gen(10, rand.Float64)
func Gen[T any](sz int, generator func() T) []T {
	result := make([]T, sz)
	for i := range result {
		result[i] = generator()
	}
	return result
}

// GenIndex builds a new slice of a given size.
// The generator function is called for each index in the new slice.
func GenIndex[T any](sz int, generator func(i int) T) []T {
	result := make([]T, sz)
	for i := range result {
		result[i] = generator(i)
	}
	return result
}

// Copy returns a shallow copy of the given slice.
func Copy[T any](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}

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

// Reverse reverses the elements of s in-place and returns s again for easy chaining.
func Reverse[T any](s []T) []T {
	end := len(s) / 2
	for i := 0; i < end; i++ {
		swapIdx := len(s) - 1 - i
		s[i], s[swapIdx] = s[swapIdx], s[i]
	}
	return s
}

// Shuffle pseudo-randomizes the elements in the slice in-place.
// Returns the slice again for easy chaining.
func Shuffle[T any](s []T) []T {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}

// First returns the first element of s or the zero value of T
func First[T any](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[0]
}

// Last returns the last element of s or the zero value of T
func Last[T any](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[len(s)-1]
}

// Delete removes all slice entries where the predicate returns true.
// The input slice is changed in-place.
func Delete[T any](s []T, shouldDelete func(T) bool) []T {
	// One pass deletion.
	// We put entries to keep to the left, and entries to delete to the right,
	// while preserving the ordering of the kept entries.
	highestKeepIdx := 0
	for idx := range s {
		if !shouldDelete(s[idx]) {
			if idx != highestKeepIdx {
				s[highestKeepIdx], s[idx] = s[idx], s[highestKeepIdx]
			}
			highestKeepIdx++
		}
	}

	// trim the columns
	return s[:highestKeepIdx]
}

// Trim removes any zero-value entries from the beginning and end of a given slice.
// The input slice is unchanged and a sub-slice is returned.
func Trim[T comparable](s []T) []T {
	if s == nil {
		return nil
	}

	var (
		zero       T
		start, end int
	)
	for start = 0; start < len(s) && s[start] == zero; start++ {
	}

	for end = len(s); end > start && s[end-1] == zero; end-- {
	}

	return s[start:end]
}
