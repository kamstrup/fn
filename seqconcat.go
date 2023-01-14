package fn

type concatSeq[T any] struct {
	head Seq[T]
	tail Seq[Seq[T]]
}

func FlattenOf[T any](seqs Seq[Seq[T]]) Seq[T] {
	return concatSeq[T]{
		head: nil,
		tail: seqs,
	}
}

func ConcatOf[T any](seqs ...Seq[T]) Seq[T] {
	return concatSeq[T]{
		head: nil,
		tail: ArrayOf(seqs),
	}
}

func (c concatSeq[T]) ForEach(f Func1[T]) Seq[T] {
	if c.head != nil {
		c.ForEach(f)
	}
	if c.tail != nil {
		c.tail.ForEach(func(seq Seq[T]) {
			seq.ForEach(f)
		})
	}

	return SeqEmpty[T]()
}

func (c concatSeq[T]) ForEachIndex(f Func2[int, T]) Seq[T] {
	i := 0
	c.ForEach(func(t T) {
		f(i, t)
		i++
	})

	return SeqEmpty[T]()
}

func (c concatSeq[T]) Len() (int, bool) {
	// fx. if c.tail is an Array we can do a stateless check to see if we can calculate a total length
	tailArr, tailIsArray := c.tail.(Array[Seq[T]])
	if !tailIsArray {
		return LenUnknown, false
	}
	sz := 0
	if c.head != nil {
		if l, ok := c.head.Len(); ok {
			sz += l
		} else {
			return l, false // head is infinite or unknown
		}
	}

	for _, tailSeq := range tailArr {
		if l, ok := tailSeq.Len(); ok {
			sz += l
		} else {
			return l, false // tailSeq is infinite or unknown
		}
	}

	return sz, true
}

func (c concatSeq[T]) Array() Array[T] {
	var buf []T
	sz, hasLen := c.Len()
	if hasLen {
		buf = make([]T, 0, sz)
	}

	buf = Into(buf, Append[T], c.seq())
	return buf
}

func (c concatSeq[T]) Take(n int) (Array[T], Seq[T]) {
	if n <= 0 {
		return []T{}, c
	}
	var (
		arr      Array[T]
		headTail Seq[T]
	)

	// First see if we have enough in c.head
	if c.head != nil {
		arr, headTail = c.head.Take(n)
		if len(arr) == n {
			return arr, concatSeq[T]{
				head: headTail,
				tail: c.tail,
			}
		}
		c.head = nil // head depleted (it returned < n)
	}

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

		headArr, headTail = headOpt.val.Take(n - len(arr))
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
	var arr []T
	for fst, tail := c.First(); fst.Ok(); fst, tail = tail.First() {
		if !pred(fst.val) {
			return arr, ConcatOf(SingletOf(fst.val), tail)
		}
		arr = append(arr, fst.val)
	}
	return arr, SeqEmpty[T]()
}

func (c concatSeq[T]) Skip(n int) Seq[T] {
	var (
		i    = 0
		fst  Opt[T]
		tail Seq[T]
	)
	for fst, tail = c.First(); fst.Ok() && i < n; fst, tail = tail.First() {
		i++
	}
	return tail
}

func (c concatSeq[T]) Where(pred Predicate[T]) Seq[T] {
	return whereSeq[T]{
		seq:  c,
		pred: pred,
	}
}

func (c concatSeq[T]) While(pred Predicate[T]) Seq[T] {
	return whileSeq[T]{
		seq:  c,
		pred: pred,
	}
}

func (c concatSeq[T]) First() (Opt[T], Seq[T]) {
	var (
		fst      Opt[T]
		headTail Seq[T]
	)

	// First see if we have enough in c.head
	if c.head != nil {
		fst, headTail = c.head.First()
		if fst.Ok() {
			return fst, concatSeq[T]{
				head: headTail,
				tail: c.tail,
			}
		}
		c.head = nil // head depleted
	}

	for {
		// If we get here, c.head is depleted, and we still need an element,
		// so check if we have a new head in c.tail

		var headOpt Opt[Seq[T]]
		headOpt, c.tail = c.tail.First()
		if headOpt.Empty() {
			// No new head, we took the last elements in initial head.First()
			return OptEmpty[T](), SeqEmpty[T]()
		}

		// We have anew head
		c.head = headOpt.val
		fst, headTail = c.head.First()
		if fst.Ok() {
			return fst, concatSeq[T]{
				head: headTail,
				tail: c.tail,
			}
		}

		// New head was empty, go again
		c.head = nil
	}
}

func (c concatSeq[T]) Map(shaper FuncMap[T, T]) Seq[T] {
	return mappedSeq[T, T]{
		f:   shaper,
		seq: c,
	}
}

func (c concatSeq[T]) Error() error {
	if c.head != nil {
		return Error(c.head)
	}
	return Error(c.tail)
}

// seq is just a cast helper, to make the Go compiler happy
func (c concatSeq[T]) seq() Seq[T] {
	return c
}
