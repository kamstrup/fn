package seq

// Tuple represents a pair of values.
// They normally show up when using ZipOf() or iterating over a Map[X,Y].
type Tuple[X comparable, Y any] struct {
	x X
	y Y
}

func TupleOf[X comparable, Y any](x X, y Y) Tuple[X, Y] {
	return Tuple[X, Y]{x, y}
}

// TupleKey returns the key of a tuple.
// A simpler variant of the method expression (seq.Tuple[k,V]).Key
func TupleKey[X comparable, Y any](t Tuple[X, Y]) X {
	return t.x
}

// TupleValue returns the value of a tuple.
// A simpler variant of the method expression (seq.Tuple[k,V]).Value
func TupleValue[X comparable, Y any](t Tuple[X, Y]) Y {
	return t.y
}

// X is an alias for Key
func (t Tuple[X, Y]) X() X {
	return t.x
}

// Key returns the first element in the tuple
func (t Tuple[X, Y]) Key() X {
	return t.x
}

// Y is an alias for Value
func (t Tuple[X, Y]) Y() Y {
	return t.y
}

// Value returns the second member of the tuple
func (t Tuple[X, Y]) Value() Y {
	return t.y
}

func (t Tuple[X, Y]) Equals(o Tuple[X, Y]) bool {
	// Weird && construction to avoid potential interface{} alloc.
	// Tuple.Y is not comparable so we have to bend it a bit
	if t.x == o.x {
		var ty, oy any = t.y, o.y
		return ty == oy
	}
	return false
}
