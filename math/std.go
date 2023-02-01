package fnmath

import "github.com/kamstrup/fn/constraints"

// Sum is a seq.FuncCollect and a seq.FuncUpdate, for use with seq.Reduce or seq.UpdateSlice, that sums up the elements it sees.
func Sum[T constraints.Arithmetic](into, t T) T {
	return into + t
}

// Max is a seq.FuncCollect and a seq.FuncUpdate, for use with seq.Reduce or seq.UpdateMap, that returns the maximal element.
func Max[T constraints.Ordered](s, t T) T {
	if s > t {
		return s
	}
	return t
}

// Min is a seq.FuncCollect and a seq.FuncUpdate, for use with seq.Reduce or seq.UpdateMap, that returns the minimal element.
func Min[T constraints.Ordered](s, t T) T {
	if s < t {
		return s
	}
	return t
}
