package fx

// GroupBy returns a map where elements with the same key are collected together in a slice.
func GroupBy[K comparable, V any](values []V, key func(V) K) map[K][]V {
	groups := make(map[K][]V)
	for _, v := range values {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}
