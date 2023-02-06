package slice

// Equal returns true the two slices contains the same elements.
// Empty slices and nil slices are considered equal.
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v1 := range s1 {
		if v1 != s2[i] {
			return false
		}
	}

	return true
}

// HasPrefix returns true if the slice s starts with the given prefix.
// If the prefix is empty (including nil) this function always returns true.
func HasPrefix[T comparable](s, prefix []T) bool {
	if len(prefix) > len(s) {
		return false
	}

	for i, p := range prefix {
		if p != s[i] {
			return false
		}
	}

	return true
}

// HasSuffix returns true if the slice s ends with the given suffix.
// If the suffix is empty (including nil) this function always returns true.
func HasSuffix[T comparable](s, suffix []T) bool {
	if len(suffix) > len(s) {
		return false
	}

	offset := len(s) - len(suffix)

	for i, sf := range suffix {
		if sf != s[offset+i] {
			return false
		}
	}

	return true
}
