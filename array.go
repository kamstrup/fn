package fn

import "sort"

type Array[T any] struct {
	vals []T
}

func SeqEmpty[T any]() Seq[T] {
	return Array[T]{}
}

func ArrayOf[T any](tt []T) Array[T] {
	return Array[T]{tt}
}

func ArrayOfArgs[T any](tt ...T) Array[T] {
	return Array[T]{tt}
}

// Seq is a helper function for letting the Go compiler understand that Array[T] implements Seq[T]
func (a Array[T]) Seq() Seq[T] {
	return a
}

func (a Array[T]) ForEach(f Func1[T]) {
	for _, v := range a.vals {
		f(v)
	}
}

func (a Array[T]) ForEachIndex(f Func2[int, T]) {
	for i, v := range a.vals {
		f(i, v)
	}
}

func (a Array[T]) Len() int {
	return len(a.vals)
}

func (a Array[T]) Array() Array[T] {
	return a
}

func (a Array[T]) Take(n int) (Array[T], Seq[T]) {
	if a.Len() <= n {
		return a, SeqEmpty[T]()
	}
	return Array[T]{vals: a.vals[:n]}, Array[T]{vals: a.vals[n:]}
}

func (a Array[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	for i, v := range a.vals {
		if !pred(v) {
			return ArrayOf(a.vals[:i]), ArrayOf(a.vals[i:])
		}
	}
	return a, SeqEmpty[T]()
}

func (a Array[T]) Skip(n int) Seq[T] {
	if a.Len() <= n {
		return SeqEmpty[T]()
	}
	return Array[T]{vals: a.vals[n:]}
}

func (a Array[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Array[T]) First() (Opt[T], Seq[T]) {
	if len(a.vals) == 0 {
		return OptEmpty[T](), a
	}
	return OptOf(a.vals[0]), ArrayOf(a.vals[1:])
}

// Sort is special for Array Seqs since it is done in place.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Array[T]) Sort(less FuncLess[T]) Array[T] {
	sort.Slice(a.vals, func(i, j int) bool {
		return less(a.vals[i], a.vals[j])
	})
	return a
}

// Reverse is special for Array Seqs since it is done in place.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Array[T]) Reverse() Seq[T] {
	end := len(a.vals) / 2
	for i := 0; i < end; i++ {
		swapIdx := len(a.vals) - 1 - i
		a.vals[i], a.vals[swapIdx] = a.vals[swapIdx], a.vals[i]
	}
	return a
}

// Slice provides raw access to the underlying data of this Array
func (a Array[T]) Slice() []T {
	return a.vals
}
