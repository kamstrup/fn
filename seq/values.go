package seq

import "github.com/kamstrup/fn/opt"

type valuesSeq[T any] struct {
	seq Seq[opt.Opt[T]]
}

// ValuesOf returns a seq that lazily converts opts into their wrapped values, stopping at the first error.
func ValuesOf[T any](opts Seq[opt.Opt[T]]) Seq[T] {
	return valuesSeq[T]{seq: opts}
}

func (v valuesSeq[T]) ForEach(f Func1[T]) Seq[T] {
	var (
		val     T
		loopErr error
	)
	res := v.seq.ForEach(func(o opt.Opt[T]) {
		if loopErr == nil { // only do something if we didn't already see an error
			val, loopErr = o.Return()
			if loopErr == nil {
				f(val)
			}
		} // else: skip remaining items
	})

	if resErr := Error(res); resErr != nil {
		return ErrorOf[T](resErr)
	} else if loopErr != nil {
		return ErrorOf[T](loopErr)
	}
	return Empty[T]()
}

func (v valuesSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	i := 0
	return v.ForEach(func(t T) {
		f(i, t)
		i++
	})
}

func (v valuesSeq[T]) Len() (int, bool) {
	return v.seq.Len() // FIXME: length could be shorter if there is an error opt
}

func (v valuesSeq[T]) ToSlice() Slice[T] {
	var arr []T
	if sz, ok := v.Len(); ok {
		arr = make([]T, 0, sz)
	}
	v.ForEach(func(t T) {
		arr = append(arr, t)
	})
	return arr
}

func (v valuesSeq[T]) Limit(n int) Seq[T] {
	return LimitOf[T](v, n)
}

func (v valuesSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	i := 0
	return v.TakeWhile(func(t T) bool {
		if i < n {
			i++
			return true
		}
		return false
	})
}

func (v valuesSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	var (
		arr  []T
		head opt.Opt[T]
		tail Seq[T]
	)
	for head, tail = v.First(); ; head, tail = tail.First() {
		t, err := head.Return()
		if err != nil {
			return arr, ErrorOf[T](err)
		} else if pred(t) {
			arr = append(arr, t)
		} else {
			return arr, PrependOf(t, tail)
		}
	}
	return arr, tail
}

func (v valuesSeq[T]) Skip(n int) Seq[T] {
	var (
		head opt.Opt[T]
		tail Seq[T]
	)
	for head, tail = v.First(); n > 0; head, tail = tail.First() {
		n--
		if err := head.Error(); err != nil {
			return ErrorOf[T](err)
		}
	}
	return tail
}

func (v valuesSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return WhereOf[T](v, pred)
}

func (v valuesSeq[T]) While(pred Predicate[T]) Seq[T] {
	return WhileOf[T](v, pred)
}

func (v valuesSeq[T]) First() (opt.Opt[T], Seq[T]) {
	fst, tail := v.seq.First() // note: fst is an Opt[Opt[T]]
	optT, err := fst.Return()
	if err != nil {
		return opt.ErrorOf[T](err), ErrorOf[T](err)
	}
	return optT, valuesSeq[T]{seq: tail}
}

func (v valuesSeq[T]) Map(funcMap FuncMap[T, T]) Seq[T] {
	return MappingOf[T, T](v, funcMap)
}
