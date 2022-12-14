package fn

type Tuple[X comparable, Y any] struct {
	x X
	y Y
}

func TupleOf[X comparable, Y any](x X, y Y) Tuple[X, Y] {
	return Tuple[X, Y]{x, y}
}

func TupleKey[X comparable, Y any](t Tuple[X, Y]) X {
	return t.x
}

func TupleValue[X comparable, Y any](t Tuple[X, Y]) Y {
	return t.y
}

func (t Tuple[X, Y]) X() X {
	return t.x
}

func (t Tuple[X, Y]) Key() X {
	return t.x
}

func (t Tuple[X, Y]) Y() Y {
	return t.y
}

func (t Tuple[X, Y]) Value() Y {
	return t.y
}
