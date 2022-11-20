package fn

func lessThan[N Ordered](n N) Predicate[N] {
	return func(i N) bool {
		return i < n
	}
}

func greaterThan[N Ordered](n N) Predicate[N] {
	return func(i N) bool {
		return i > n
	}
}

// RangeOf returns a Seq that counts from one number to another.
// It can count both up or down.
func RangeOf(from, to int) Seq[int] {
	if from == to {
		return SeqEmpty[int]()
	} else if from < to {
		return SourceOf(NumbersFrom(from)).While(lessThan(to))
	} else {
		return SourceOf(NumbersLowerThan(from)).While(greaterThan(to))
	}

}
