package fn

type rangeSeq[N Integer] struct {
	from N
	to   N
	step N // we have to work with a positive step in order to support unsigned ints
}

// RangeStepOf returns a Seq that counts from one number to another, in increments of some given step.
// It can count both up or down. Range Seqs have a well-defined length.
func RangeStepOf[N Integer](from, to, step N) Seq[N] {
	if step == 0 {
		panic("range step must be non-zero")
	}
	if from == to {
		return SeqEmpty[N]()
	}

	// We must work with a positive step, in order to support unsigned numbers
	if step < 0 {
		step = -step
	}

	return rangeSeq[N]{
		from: from,
		to:   to,
		step: step,
	}
}

// RangeOf returns a Seq that counts from one number to another.
// It can count both up or down. Range Seqs have a well-defined length.
func RangeOf[N Integer](from, to N) Seq[N] {
	return RangeStepOf(from, to, 1)
}

func (r rangeSeq[N]) ForEach(f Func1[N]) Seq[N] {
	if r.to >= r.from {
		for i := r.from; i < r.to; i += r.step {
			f(i)
		}
	} else {
		for i := r.from; i > r.to && i <= r.from; i -= r.step {
			f(i)
		}
	}

	return SeqEmpty[N]()
}

func (r rangeSeq[N]) ForEachIndex(f Func2[int, N]) Seq[N] {
	i := 0
	if r.to >= r.from {
		for j := r.from; j < r.to; j += r.step {
			f(i, j)
			i++
		}
	} else {
		for j := r.from; j > r.to && j <= r.from; j -= r.step {
			f(i, j)
			i++
		}
	}

	return SeqEmpty[N]()
}

func (r rangeSeq[N]) Len() (int, bool) {
	var rng N
	if r.to >= r.from {
		rng = r.to - r.from
	} else {
		rng = r.from - r.to
	}

	remainder := 0
	if rng%r.step != 0 {
		remainder = 1
	}
	sz := int(rng / r.step)
	if sz >= 0 {
		return sz + remainder, true
	}
	return -sz + remainder, true
}

func (r rangeSeq[N]) Array() Array[N] {
	sz, _ := r.Len()
	arr := make([]N, sz)
	i := 0
	if r.to >= r.from {
		for j := r.from; j < r.to; j += r.step {
			arr[i] = j
			i++
		}
	} else {
		for j := r.from; j > r.to && j <= r.from; j -= r.step {
			arr[i] = j
			i++
		}
	}

	return arr
}

func (r rangeSeq[N]) Take(n int) (Array[N], Seq[N]) {
	sz, _ := r.Len()
	if n >= sz {
		return r.Array(), SeqEmpty[N]()
	}

	// n < sz
	arr := make([]N, n)
	if r.to > r.from {
		for i := 0; i < n; i++ {
			arr[i] = r.from + r.step*N(i)
		}
		return arr, RangeStepOf(r.from+N(n)*r.step, r.to, r.step)
	} else {
		for i := 0; i < n; i++ {
			arr[i] = r.from - r.step*N(i)
		}
		return arr, RangeStepOf(r.from-N(n)*r.step, r.to, r.step)
	}
}

func (r rangeSeq[N]) TakeWhile(pred Predicate[N]) (Array[N], Seq[N]) {
	sz, _ := r.Len()
	var arr []N
	if r.to > r.from {
		for i := 0; i < sz; i++ {
			val := r.from + r.step*N(i)
			if pred(val) {
				arr = append(arr, val)
			} else {
				return arr, RangeStepOf(r.from+N(i)*r.step, r.to, r.step)
			}
		}
	} else {
		for i := 0; i < sz; i++ {
			val := r.from - r.step*N(i)
			if pred(val) {
				arr = append(arr, val)
			} else {
				return arr, RangeStepOf(r.from+N(i)*r.step, r.to, r.step)
			}
		}
	}

	return arr, SeqEmpty[N]()
}

func (r rangeSeq[N]) Skip(n int) Seq[N] {
	if n < 0 {
		panic("must skip >= 0 elements")
	}

	if sz, _ := r.Len(); n >= sz {
		return SeqEmpty[N]()
	}

	if r.to > r.from {
		return RangeStepOf(r.from+N(n)*r.step, r.to, r.step)
	}
	return RangeStepOf(r.from-N(n)*r.step, r.to, r.step)

}

func (r rangeSeq[N]) Where(pred Predicate[N]) Seq[N] {
	return whereSeq[N]{
		seq:  r,
		pred: pred,
	}
}

func (r rangeSeq[N]) While(pred Predicate[N]) Seq[N] {
	return whileSeq[N]{
		seq:  r,
		pred: pred,
	}
}

func (r rangeSeq[N]) First() (Opt[N], Seq[N]) {
	if r.from <= r.to {
		if v := r.from + r.step; v < r.to {
			return OptOf(r.from), RangeStepOf(v, r.to, r.step)
		}
	} else {
		if v := r.from - r.step; v > r.to && v < r.from {
			return OptOf(r.from), RangeStepOf(v, r.to, r.step)
		}
	}

	return OptOf(r.from), SeqEmpty[N]()
}

func (r rangeSeq[N]) Map(shaper FuncMap[N, N]) Seq[N] {
	return mappedSeq[N, N]{
		f:   shaper,
		seq: r,
	}
}
