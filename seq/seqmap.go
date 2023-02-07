package seq

import "github.com/kamstrup/fn/opt"

// MappingOf creates a new Seq that lazily converts the values, via a FuncMap,
// into another Seq. The returned Seq has the same Seq.Len() as the input Seq.
// If you are looking for ways to create a Seq from a go map[K]V please
// look at MapOf() or SetOf().
//
// If the mapping function is some kind of heavy operation or requires IO,
// consider using the parallelized version of MappingOf called Go.
//
// If the mapping operation can result in an error consider wrapping the return
// type in an opt. This can be done easily with opt.Mapper. For example,
// parsing strings as ints and discarding any errors could look like:
//
//	  strInts := seq.SliceOfArgs("1", "two", "3")
//		 ints := seq.MappingOf(strInts, opt.Mapper(strconv.Atoi)).
//			         Where(opt.Ok[int])
func MappingOf[S, T any](seq Seq[S], fm FuncMap[S, T]) Seq[T] {
	return mappedSeq[S, T]{
		f:   fm,
		seq: seq,
	}
}

type mappedSeq[S, T any] struct {
	f   FuncMap[S, T]
	seq Seq[S]
}

func (m mappedSeq[S, T]) ForEach(f Func1[T]) Seq[T] {
	res := m.seq.ForEach(func(s S) {
		t := m.f(s)
		f(t)
	})

	if err := Error(res); err != nil {
		return ErrorOf[T](err)
	}
	return Empty[T]()
}

func (m mappedSeq[S, T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	res := m.seq.ForEachIndex(func(i int, s S) {
		t := m.f(s)
		f(i, t)
	})

	if err := Error(res); err != nil {
		return ErrorOf[T](err)
	}
	return Empty[T]()
}

func (m mappedSeq[S, T]) Len() (int, bool) {
	return m.seq.Len()
}

func (m mappedSeq[S, T]) ToSlice() Slice[T] {
	if sz, ok := m.seq.Len(); ok {
		arr := make([]T, sz)
		m.ForEachIndex(func(i int, t T) {
			arr[i] = t
		})
		return arr
	} else {
		var arr []T
		m.ForEach(func(t T) {
			arr = append(arr, t)
		})
		return arr
	}
}

func (m mappedSeq[S, T]) Take(n int) (Slice[T], Seq[T]) {
	var (
		head Slice[S]
		tail Seq[S]
	)
	// Note: we are not calling m.f on the skipped elements
	head, tail = m.seq.Take(n)
	return MappingOf[S, T](head, m.f).ToSlice(), MappingOf(tail, m.f)
}

func (m mappedSeq[S, T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	// TODO: does not really need to alloc a slice, if we had a "pulling seq"
	// FIXME: does double alloc with TakeWhile+append, maybe just While+Do?
	var arr []T
	_, tail := m.seq.TakeWhile(func(s S) bool {
		t := m.f(s)
		if pred(t) {
			arr = append(arr, t)
			return true
		}
		return false
	})
	return arr, MappingOf(tail, m.f)
}

func (m mappedSeq[S, T]) Skip(n int) Seq[T] {
	tail := m.seq.Skip(n)
	return MappingOf(tail, m.f)
}
func (m mappedSeq[S, T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  m,
		pred: pred,
	}
}

func (m mappedSeq[S, T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  m,
		pred: pred,
	}
}

func (m mappedSeq[S, T]) First() (opt.Opt[T], Seq[T]) {
	s, tail := m.seq.First()
	return opt.Map(s, m.f), MappingOf(tail, m.f)
}

func (m mappedSeq[S, T]) Map(shaper FuncMap[T, T]) Seq[T] {
	return mappedSeq[T, T]{
		f:   shaper,
		seq: m,
	}
}

func (m mappedSeq[S, T]) Error() error {
	return Error(m.seq)
}
