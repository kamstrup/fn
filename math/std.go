package fnmath

import "github.com/kamstrup/fn/constraints"

// Sum is a fn.FuncCollect and a FuncUpdate, for use with fn.Into or fn.UpdateAt, that sums up the elements it sees.
func Sum[T constraints.Arithmetic](into, t T) T {
	return into + t
}

// Max is a fn.FuncCollect and a fn.FuncUpdate, for use with fn.Into or fn.UpdateAssoc, that returns the maximal element.
func Max[T constraints.Ordered](s, t T) T {
	if s > t {
		return s
	}
	return t
}

// Min is a fn.FuncCollect and a fn.FuncUpdate, for use with fn.Into or fn.UpdateAssoc, that returns the minimal element.
func Min[T constraints.Ordered](s, t T) T {
	if s < t {
		return s
	}
	return t
}
