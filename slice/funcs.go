package slice

import "math/bits"

// GroupBy returns a map where elements with the same key are collected together in a slice.
func GroupBy[K comparable, V any](values []V, key func(V) K) map[K][]V {
	groups := make(map[K][]V)
	for _, v := range values {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}

// Uniq returns a slice with all the unique values from the input slice, in the order they appear.
func Uniq[T comparable](slice []T) []T {
	sz := simpleLog(len(slice))
	dict := make(map[T]struct{}, sz)
	var result []T
	for _, t := range slice {
		if _, seen := dict[t]; !seen {
			dict[t] = struct{}{}
			result = append(result, t)
		}
	}

	return result
}

// simpleLog calculates a simplified log2 of an int.
// We use this to estimate sizes of various dictionaries,
// since the number of unique words in a corpus normally
// grows logarithmically in the number of words.
func simpleLog(corpusSize int) int {
	return bits.Len(uint(corpusSize))
}

// Reduce collects a slice of elements E into a result of type R.
// This operation is also known as "fold" in other frameworks.
//
// It is a simplified form of seq.Reduce found in the core package, and works with
// standard seq.FuncCollect functions like seq.MakeString, seq.MakeBytes, seq.Count,
// seq.MakeSet, fnmath.Min, fnmath.Max, fnmath.Sum etc.
//
// The first argument "result" is the initial state. The second is the collector
// function. The signature of the collector works like standard append(). The last
// argument is the slice to collect values from.
//
// # Example:
//
// Summing the elements in and integer slice could look like:
//
//	vals := []int{1, 2, 3}
//	sum := slice.Reduce(func(res, e int) int {
//	    return res + e
//	}, 0, vals)
//
//	// sum is now: 6
func Reduce[E any, R any](collector func(R, E) R, result R, from []E) R {
	for _, e := range from {
		result = collector(result, e)
	}
	return result
}
