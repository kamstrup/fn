package slice

// Keys returns a slice with all keys from a map.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
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

// FromMap builds a slice []T from a map[K]V.
// To build a map from a slice use ToMap.
// Please bear in mind that go maps are unordered and the resulting slice will reflect that.
func FromMap[K comparable, V, T any](m map[K]V, f func(K, V) T) []T {
	slice := make([]T, len(m))
	i := 0
	for k, v := range m {
		slice[i] = f(k, v)
		i++
	}
	return slice
}

// ToMap builds a map[K]V from a slice.
// To build a slice from a map use FromMap.
// The mnemonic for the function name is that the last word is what you operate on.
func ToMap[K comparable, V, T any](slice []T, f func(T) (K, V)) map[K]V {
	m := make(map[K]V, len(slice))
	for _, t := range slice {
		k, v := f(t)
		m[k] = v
	}
	return m
}
