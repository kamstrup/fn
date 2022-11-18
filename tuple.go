package fn

type Tuple[X, Y comparable] struct {
	x X
	y Y
}

func TupleOf[X, Y comparable](x X, y Y) Tuple[X, Y] {
	return Tuple[X, Y]{x, y}
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
