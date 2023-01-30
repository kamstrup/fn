package fx

// Into collects a slice of elements E into a result of type R.
// This operation is also known as "fold" or "reduce" in other frameworks.
//
// It is a simplified form of seq.Into found in the core package, and works with
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
//	 vals := []int{1, 2, 3}
//		sum := Into(0, func(res, e int) int {
//			return res + e
//		}, vals)
//
//	 // sum is now: 6
func Into[E any, R any](result R, collect func(R, E) R, from []E) R {
	for _, e := range from {
		result = collect(result, e)
	}
	return result
}
