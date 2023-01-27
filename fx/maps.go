package fx

// Keys returns a slice with all keys from a map.
func Keys[K comparable, V any](m map[K]V) []V {
	keys := make([]V, 0, len(m))
	for _, v := range m {
		keys = append(keys, v)
	}
	return keys
}

// Values returns a slice with all keys from a map.
func Values[K comparable, V any](m map[K]V) []V {
	vals := make([]V, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// MapAssoc converts a `map[K1]V1` to a `map[K2]V2`.
func MapAssoc[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2)) map[K2]V2 {
	m2 := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := f(k, v)
		m2[k2] = v2
	}
	return m2
}

// SliceAssoc builds a slice []T from a map[K]V.
// To build a map from a slice use AssocSlice.
// The mnemonic for the function name is that the last word is what you operate on.
func SliceAssoc[K comparable, V, T any](m map[K]V, f func(K, V) T) []T {
	slice := make([]T, len(m))
	i := 0
	for k, v := range m {
		slice[i] = f(k, v)
		i++
	}
	return slice
}
