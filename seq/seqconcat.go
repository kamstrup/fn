package seq

import "github.com/kamstrup/fn/opt"

type concatSeq[T any] struct {
	head Seq[T]
	tail Seq[Seq[T]]
}

// FlattenOf returns a lazy Seq that steps through elements in a collection of seqs as though it was one big seq.
// See also ConcatOf and PrependOf.
func FlattenOf[T any](seqs Seq[Seq[T]]) Seq[T] {
	return concatSeq[T]{
		head: nil,
		tail: seqs,
	}
}

// ConcatOf wraps a collection of seqs as one contiguous lazy seq.
// See also FlattenOf and PrependOf.
func ConcatOf[T any](seqs ...Seq[T]) Seq[T] {
	return concatSeq[T]{
		head: nil,
		tail: SliceOf(seqs),
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
	return LenUnknown, false
	// fx. if c.tail is an Slice we can do a stateless check to see if we can calculate a total length
	/*tailArr, tailIsArray := c.tail.(Slice[Seq[T]])
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

	return sz, true*/
}

func (c concatSeq[T]) Values() Slice[T] {
	var buf []T
	sz, hasLen := c.Len()
	if hasLen {
		buf = make([]T, 0, sz)
	}

	buf = Reduce(Append[T], buf, c.seq()).Or(nil) // careful: errors silently dropped
	return buf
}

func (c concatSeq[T]) Take(n int) (Slice[T], Seq[T]) {
	if n <= 0 {
		return []T{}, c
	}
	var (
		arr      Slice[T]
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
			headOpt opt.Opt[Seq[T]]
			headArr Slice[T]
		)
		headOpt, c.tail = c.tail.First()
		headSeq, headErr := headOpt.Return()
		if headErr != nil {
			// No new head, we took the last elements in initial head.Take(n)
			return arr, ErrorOf[T](headErr)
		}

		headArr, headTail = headSeq.Take(n - len(arr))
		arr = append(arr, headArr...)
		c.head = headTail

		if len(arr) == n {
			return arr, c
		}

		// c.head depleted, go again
		c.head = nil
	}
}

func (c concatSeq[T]) TakeWhile(pred Predicate[T]) (Slice[T], Seq[T]) {
	var (
		arr  []T
		fst  opt.Opt[T]
		tail Seq[T]
	)
	for fst, tail = c.First(); fst.Ok(); fst, tail = tail.First() {
		val := fst.Must()
		if !pred(val) {
			return arr, PrependOf(val, tail)
		}
		arr = append(arr, val)
	}
	return arr, ErrorOf[T](fst.Error())
}

func (c concatSeq[T]) Skip(n int) Seq[T] {
	var (
		i    = 0
		fst  opt.Opt[T]
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

func (c concatSeq[T]) First() (opt.Opt[T], Seq[T]) {
	var (
		fst      opt.Opt[T]
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

		var headOpt opt.Opt[Seq[T]]
		headOpt, c.tail = c.tail.First()

		headVal, headErr := headOpt.Return()
		if headErr != nil {
			// No new head, we took the last elements in initial head.First()
			return opt.ErrorOf[T](headErr), ErrorOf[T](headErr)
		}

		// We have a new head
		c.head = headVal
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
