package fn

import "sort"

type Array[T any] []T

func SeqEmpty[T any]() Seq[T] {
	return Array[T](nil)
}

func ArrayOf[T any](tt []T) Array[T] {
	return tt
}

func ArrayOfArgs[T any](tt ...T) Array[T] {
	return tt
}

// Seq is a helper function for letting the Go compiler understand that Array[T] implements Seq[T]
func (a Array[T]) Seq() Seq[T] {
	return a
}

func (a Array[T]) ForEach(f Func1[T]) {
	for _, v := range a {
		f(v)
	}
}

func (a Array[T]) ForEachIndex(f Func2[int, T]) {
	for i, v := range a {
		f(i, v)
	}
}

func (a Array[T]) Len() int {
	return len(a)
}

func (a Array[T]) Array() Array[T] {
	return a
}

func (a Array[T]) Take(n int) (Array[T], Seq[T]) {
	if a.Len() <= n {
		return a, SeqEmpty[T]()
	}
	return a[:n], a[n:]
}

func (a Array[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	for i, v := range a {
		if !pred(v) {
			return a[:i], a[i:]
		}
	}
	return a, SeqEmpty[T]()
}

func (a Array[T]) Skip(n int) Seq[T] {
	if a.Len() <= n {
		return SeqEmpty[T]()
	}
	return a[n:]
}

func (a Array[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Array[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  a,
		pred: pred,
	}
}

func (a Array[T]) First() (Opt[T], Seq[T]) {
	if len(a) == 0 {
		return OptEmpty[T](), a
	}
	return OptOf(a[0]), ArrayOf(a[1:])
}

// Sort is special for Array Seqs since it is done in place.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
//
// Typical less-functions to use are OrderAsc, OrderDesc, OrderTupleAsc, and OrderTupleDesc.
func (a Array[T]) Sort(less FuncLess[T]) Array[T] {
	sort.Slice(a, func(i, j int) bool {
		return less(a[i], a[j])
	})
	return a
}

// Reverse is special for Array Seqs since it is done in place.
// Generally functions and methods in the fn() library leaves all data structures immutable,
// but this is an exception. Caveat Emptor!
func (a Array[T]) Reverse() Seq[T] {
	end := len(a) / 2
	for i := 0; i < end; i++ {
		swapIdx := len(a) - 1 - i
		a[i], a[swapIdx] = a[swapIdx], a[i]
	}
	return a
}

// AsSlice is a chainable method for casting the Array into a []T
func (a Array[T]) AsSlice() []T {
	return a
}
