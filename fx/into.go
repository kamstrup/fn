package fx

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
//	 vals := []int{1, 2, 3}
//		sum := Reduce(0, func(res, e int) int {
//			return res + e
//		}, vals)
//
//	 // sum is now: 6
func Reduce[E any, R any](collector func(R, E) R, result R, from []E) R {
	for _, e := range from {
		result = collector(result, e)
	}
	return result
}
