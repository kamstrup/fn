package fnmath

import (
	"github.com/kamstrup/fn/constraints"
)

// Stats is the result structure used by MakeStats
type Stats[N constraints.Arithmetic] struct {
	Sum, Min, Max N
	Count         int
}

// MakeStats is a seq.FuncCollect that computes basic numeric properties of a seq of numbers.
//
// # Example
//
//	 data := seq.SliceOfArgs(-1, 0, 1, 2, 3)
//		optStats := seq.Reduce(MakeStats[int], Stats[int]{}, data)
//	 // The stats are now {Sum: 5, Min: -1, Max: 3, Count: 5}
func MakeStats[N constraints.Integer | constraints.Float](s Stats[N], n N) Stats[N] {
	if s.Count == 0 || n > s.Max {
		s.Max = n
	}
	if s.Count == 0 || n < s.Min {
		s.Min = n
	}
	s.Sum += n
	s.Count++

	return s
}
