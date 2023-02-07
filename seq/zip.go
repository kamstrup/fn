package seq

import "github.com/kamstrup/fn/opt"

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
		fx opt.Opt[X]
		fy opt.Opt[Y]
		tx = z.sx
		ty = z.sy
	)
	for {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Ok() && fy.Ok() {
			f(Tuple[X, Y]{fx.Must(), fy.Must()})
		} else {
			// we are done, at least one Seq drained
			return zipErrorTail(fx, fy)
		}
	}

	return Empty[Tuple[X, Y]]()
}

func (z zipSeq[X, Y]) ForEachIndex(f Func2[int, Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	var (
		fx opt.Opt[X]
		fy opt.Opt[Y]
		tx = z.sx
		ty = z.sy
	)
	for i := 0; ; i++ {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Ok() && fy.Ok() {
			f(i, Tuple[X, Y]{fx.Must(), fy.Must()})
		} else {
			return zipErrorTail(fx, fy) // we are done, at least one Seq drained
		}

	}

	return Empty[Tuple[X, Y]]()
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

func (z zipSeq[X, Y]) ToSlice() Slice[Tuple[X, Y]] {
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

func (z zipSeq[X, Y]) Take(n int) (Slice[Tuple[X, Y]], Seq[Tuple[X, Y]]) {
	var (
		arr []Tuple[X, Y]
		fx  opt.Opt[X]
		fy  opt.Opt[Y]
		tx  = z.sx
		ty  = z.sy
	)
	if sz, ok := z.Len(); ok {
		if sz <= n {
			// Best case
			return z.ToSlice(), Empty[Tuple[X, Y]]()
		}

		// We know we have at least n elements in both tx and ty
		arr = make([]Tuple[X, Y], n)
		for i := 0; i < n; i++ {
			fx, tx = tx.First()
			fy, ty = ty.First()
			arr[i] = Tuple[X, Y]{fx.Must(), fy.Must()}
		}
		return arr, zipSeq[X, Y]{sx: tx, sy: ty}
	}

	// Length of at least one Seq is unknown
	for i := 0; i < n; i++ {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Ok() && fy.Ok() {
			arr = append(arr, Tuple[X, Y]{fx.Must(), fy.Must()})
		} else {
			// we are done, at least one Seq drained
			return arr, zipErrorTail(fx, fy)
		}
	}

	return arr, zipSeq[X, Y]{sx: tx, sy: ty}
}

func (z zipSeq[X, Y]) TakeWhile(predicate Predicate[Tuple[X, Y]]) (Slice[Tuple[X, Y]], Seq[Tuple[X, Y]]) {
	var (
		arr []Tuple[X, Y]
		fx  opt.Opt[X]
		fy  opt.Opt[Y]
		tx  = z.sx
		ty  = z.sy
	)

	// Length of at least one Seq is unknown
	for {
		fx, tx = tx.First()
		fy, ty = ty.First()
		if fx.Empty() || fy.Empty() {
			// we are done, at least one Seq drained
			return arr, zipErrorTail(fx, fy)
		}
		tup := Tuple[X, Y]{fx.Must(), fy.Must()}
		if predicate(tup) {
			arr = append(arr, tup)
		} else {
			// pred(tup) is false, so we return.
			// We already consumed the heads of tx and ty, so we need to "put them back"
			return arr, PrependOf(tup, ZipOf(tx, ty))
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

func (z zipSeq[X, Y]) First() (opt.Opt[Tuple[X, Y]], Seq[Tuple[X, Y]]) {
	fx, tx := z.sx.First()
	fy, ty := z.sy.First()
	if errX := fx.Error(); errX != nil {
		return opt.ErrorOf[Tuple[X, Y]](errX), ErrorOf[Tuple[X, Y]](errX)
	} else if errY := fy.Error(); errY != nil {
		return opt.ErrorOf[Tuple[X, Y]](errY), ErrorOf[Tuple[X, Y]](errY)
	}

	return opt.Of(Tuple[X, Y]{fx.Must(), fy.Must()}), zipSeq[X, Y]{tx, ty}
}

func (z zipSeq[X, Y]) Map(shaper FuncMap[Tuple[X, Y], Tuple[X, Y]]) Seq[Tuple[X, Y]] {
	return mappedSeq[Tuple[X, Y], Tuple[X, Y]]{
		f:   shaper,
		seq: z,
	}
}

func (z zipSeq[X, Y]) Error() error {
	if errX := Error(z.sx); errX != nil {
		return errX
	}
	return Error(z.sy)
}

func zipErrorTail[X comparable, Y any](fx opt.Opt[X], fy opt.Opt[Y]) Seq[Tuple[X, Y]] {
	if errX := fx.Error(); errX != nil && errX != opt.ErrEmpty {
		return ErrorOf[Tuple[X, Y]](errX)
	} else if errY := fy.Error(); errY != nil && errY != opt.ErrEmpty {
		return ErrorOf[Tuple[X, Y]](errY)
	}
	return Empty[Tuple[X, Y]]()
}
