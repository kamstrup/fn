package seq_test

import (
	"reflect"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestSeqAssoc(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	as := seq.MapOf(m)

	fntesting.TestOf(t, as).LenIs(3)
	m2 := seq.Reduce(seq.MakeMap[string, int], nil, as)

	if !reflect.DeepEqual(m, m2.Must()) {
		t.Errorf("Expected %v, found %v", m, m2)
	}

	as = as.Values().Sort(seq.OrderTupleAsc[string, int])

	fntesting.TestOf(t, as).Is(
		seq.TupleOf("one", 1),
		seq.TupleOf("three", 3),
		seq.TupleOf("two", 2))
}

func TestSeqAssocWhere(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}
	as := seq.MapOf(m).Where(func(t seq.Tuple[string, int]) bool {
		return t.Key() == "one"
	})

	m2 := seq.Reduce(seq.MakeMap[string, int], nil, as)

	if !reflect.DeepEqual(map[string]int{"one": 1}, m2.Must()) {
		t.Errorf("Expected %v, found %v", m, m2)
	}
}

func TestSeqAssocSkip(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	as := seq.MapOf(m).Skip(0)
	m2 := seq.Reduce(seq.MakeMap[string, int], nil, as)
	if !reflect.DeepEqual(m, m2.Must()) {
		t.Errorf("Expected %v, found %v", m, m2)
	}

	as = seq.MapOf(m).Skip(1)
	m2 = seq.Reduce(seq.MakeMap[string, int], nil, as)
	if !(reflect.DeepEqual(map[string]int{"one": 1}, m2.Must()) || reflect.DeepEqual(map[string]int{"two": 2}, m2.Must())) {
		t.Errorf("Expected {'one':1} or {'two': 2}, found %v", m2)
	}

	as = seq.MapOf(m).Skip(123)
	m2 = seq.Reduce(seq.MakeMap[string, int], nil, as)
	if !m2.Empty() {
		t.Errorf("Expected empty map, found %v", m2.Must())
	}
}

func TestSeqAssocTake(t *testing.T) {
	t.Skip("TODO")
}

func TestSeqAssocTakeWhile(t *testing.T) {
	t.Skip("TODO")
}
