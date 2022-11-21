package fn

type concatSeq[T any] struct {
	head Seq[T]
	tail Seq[Seq[T]]
}

func ConcatOf[T any](seqs Seq[Seq[T]]) Seq[T] {
	return concatSeq[T]{
		head: nil,
		tail: seqs,
	}
}

func ConcatOfArgs[T any](seqs ...Seq[T]) Seq[T] {
	return concatSeq[T]{
		head: nil,
		tail: ArrayOf(seqs),
	}
}

func (c concatSeq[T]) ForEach(f Func1[T]) {
	if c.head != nil {
		c.ForEach(f)
	}
	if c.tail != nil {
		c.tail.ForEach(func(seq Seq[T]) {
			seq.ForEach(f)
		})
	}
}

func (c concatSeq[T]) ForEachIndex(f Func2[int, T]) {
	i := 0
	c.ForEach(func(t T) {
		f(i, t)
		i++
	})
}

func (c concatSeq[T]) Len() int {
	// TODO: There are certain cases where we can safely try calculate the total length
	// fx. if c.tail is an Array
	return LenUnknown
}

func (c concatSeq[T]) Array() Array[T] {
	// TODO: If we implement proper Len() calculation, we can pre-alloc the output array here
	arr := Into(nil, Append[T], c.seq())
	return ArrayOf(arr)
}

func (c concatSeq[T]) Take(n int) (Array[T], Seq[T]) {
	if n == 0 {
		return ArrayOf([]T{}), c
	}
	var (
		arr      Array[T]
		headTail Seq[T]
	)

	// First see if we have enough in c.head
	if c.head != nil {
		arr, headTail = c.head.Take(n)
		if arr.Len() == n {
			return arr, concatSeq[T]{
				head: headTail,
				tail: c.tail,
			}
		}
	}

	c.head = nil
	for {
		// If we get here, c.head is depleted, and arr still needs elements,
		// so check if we have a new head in c.tail

		var (
			headOpt Opt[Seq[T]]
			headArr Array[T]
		)
		headOpt, c.tail = c.tail.First()
		if headOpt.Empty() {
			// No new head, we took the last elements in initial head.Take(n)
			return arr, SeqEmpty[T]()
		}

		headArr, headTail = headOpt.val.Take(n - arr.Len())
		arr = append(arr, headArr...)
		c.head = headTail

		if len(arr) == n {
			return arr, c
		}

		// c.head depleted, go again
		c.head = nil
	}

}

func (c concatSeq[T]) TakeWhile(pred Predicate[T]) (Array[T], Seq[T]) {
	// TODO implement me
	panic("implement me")
}

func (c concatSeq[T]) Skip(n int) Seq[T] {
	// TODO implement me
	panic("implement me")
}

func (c concatSeq[T]) Where(pred Predicate[T]) Seq[T] {
	// TODO implement me
	panic("implement me")
}

func (c concatSeq[T]) While(pred Predicate[T]) Seq[T] {
	// TODO implement me
	panic("implement me")
}

func (c concatSeq[T]) First() (Opt[T], Seq[T]) {
	// TODO implement me
	panic("implement me")
}

// seq is just a cast helper, to make the Go compiler happy
func (c concatSeq[T]) seq() Seq[T] {
	return c
}
