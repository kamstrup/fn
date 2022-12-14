package fn

type zipSeq[X comparable, Y any] struct {
	sx Seq[X]
	sy Seq[Y]
}

// ZipOf creates a Seq that merges two Seqs into a series of Tuple.
// The zip Seq will stop at the shortest of the two Seqs. This protects from infinite loops
// if one of the Seqs is fx. a SourceOf().
func ZipOf[X comparable, Y any](sx Seq[X], sy Seq[Y]) Seq[Tuple[X, Y]] {
	return zipSeq[X, Y]{
		sx: sx,
		sy: sy,
	}
}

func (z zipSeq[X, Y]) ForEach(f Func1[Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	var (
		fx Opt[X]
		fy Opt[Y]
		tx = z.sx
		ty = z.sy
	)
	for {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Empty() || fy.Empty() {
			return SeqEmpty[Tuple[X, Y]]() // we are done, at least one Seq drained
		}
		f(Tuple[X, Y]{fx.val, fy.val})
	}

	return SeqEmpty[Tuple[X, Y]]()
}

func (z zipSeq[X, Y]) ForEachIndex(f Func2[int, Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	var (
		fx Opt[X]
		fy Opt[Y]
		tx = z.sx
		ty = z.sy
	)
	for i := 0; ; i++ {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Empty() || fy.Empty() {
			return SeqEmpty[Tuple[X, Y]]() // we are done, at least one Seq drained
		}
		f(i, Tuple[X, Y]{fx.val, fy.val})
	}

	return SeqEmpty[Tuple[X, Y]]()
}

func (z zipSeq[X, Y]) Len() (int, bool) {
	lx, okx := z.sx.Len()
	ly, oky := z.sy.Len()
	if okx && oky {
		// Both well-defined. Return the minimum
		if lx < ly {
			return lx, true
		}
		return ly, true
	}

	// If one is infinite, return the other.
	// This might help with allocations in certain cases.
	if lx == LenInfinite {
		return ly, oky
	} else if ly == LenInfinite {
		return lx, okx
	}

	return LenUnknown, false
}

func (z zipSeq[X, Y]) Array() Array[Tuple[X, Y]] {
	if sz, ok := z.Len(); ok {
		arr := make([]Tuple[X, Y], sz)
		z.ForEachIndex(func(i int, t Tuple[X, Y]) {
			arr[i] = t
		})
		return arr
	}

	var arr []Tuple[X, Y]
	z.ForEach(func(t Tuple[X, Y]) {
		arr = append(arr, t)
	})
	return arr
}

func (z zipSeq[X, Y]) Take(n int) (Array[Tuple[X, Y]], Seq[Tuple[X, Y]]) {
	var (
		arr []Tuple[X, Y]
		fx  Opt[X]
		fy  Opt[Y]
		tx  = z.sx
		ty  = z.sy
	)
	if sz, ok := z.Len(); ok {
		if sz <= n {
			// Best case
			return z.Array(), SeqEmpty[Tuple[X, Y]]()
		}

		// We know we have at least n elements in both tx and ty
		arr = make([]Tuple[X, Y], n)
		for i := 0; i < n; i++ {
			fx, tx = tx.First()
			fy, ty = ty.First()
			arr[i] = Tuple[X, Y]{fx.val, fy.val}
		}
		return arr, zipSeq[X, Y]{sx: tx, sy: ty}
	}

	// Length of at least one Seq is unknown
	for i := 0; i < n; i++ {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Empty() || fy.Empty() {
			// we are done, at least one Seq drained
			return arr, SeqEmpty[Tuple[X, Y]]()
		}
		arr = append(arr, Tuple[X, Y]{fx.val, fy.val})
	}

	return arr, zipSeq[X, Y]{sx: tx, sy: ty}
}

func (z zipSeq[X, Y]) TakeWhile(predicate Predicate[Tuple[X, Y]]) (Array[Tuple[X, Y]], Seq[Tuple[X, Y]]) {
	var (
		arr []Tuple[X, Y]
		fx  Opt[X]
		fy  Opt[Y]
		tx  = z.sx
		ty  = z.sy
	)

	// Length of at least one Seq is unknown
	for {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Empty() || fy.Empty() {
			// we are done, at least one Seq drained
			return arr, SeqEmpty[Tuple[X, Y]]()
		}
		tup := Tuple[X, Y]{fx.val, fy.val}
		if predicate(tup) {
			arr = append(arr, tup)
		} else {
			// pred(tup) is false, so we return.
			// We already consumed the heads of tx and ty, so we need to "put them back",
			// we do this by creating a concat() of the consumed tuple with the tail
			return arr, ConcatOf(
				SingletOf(TupleOf(fx.val, fy.val)),
				ZipOf(tx, ty))
		}
	}
}

func (z zipSeq[X, Y]) Skip(n int) Seq[Tuple[X, Y]] {
	return zipSeq[X, Y]{
		sx: z.sx.Skip(n),
		sy: z.sy.Skip(n),
	}
}

func (z zipSeq[X, Y]) Where(p Predicate[Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	return whereSeq[Tuple[X, Y]]{
		seq:  z,
		pred: p,
	}
}

func (z zipSeq[X, Y]) While(p Predicate[Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	return whileSeq[Tuple[X, Y]]{
		seq:  z,
		pred: p,
	}
}

func (z zipSeq[X, Y]) First() (Opt[Tuple[X, Y]], Seq[Tuple[X, Y]]) {
	fx, tx := z.sx.First()
	fy, ty := z.sy.First()
	if fx.Empty() || fy.Empty() {
		// we are done, at least one Seq drained
		return OptEmpty[Tuple[X, Y]](), SeqEmpty[Tuple[X, Y]]()
	}

	return OptOf(Tuple[X, Y]{fx.val, fy.val}), zipSeq[X, Y]{tx, ty}
}

func (z zipSeq[X, Y]) Map(shaper FuncMap[Tuple[X, Y], Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	return mappedSeq[Tuple[X, Y], Tuple[X, Y]]{
		f:   shaper,
		seq: z,
	}
}
