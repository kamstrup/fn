package fn_test

import (
	"reflect"
	"testing"

	"github.com/kamstrup/fn"
)

func TestSeqAssoc(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	as := fn.AssocOf(m)

	fn.SeqTest(t, as).LenIs(3)
	m2 := fn.Into(nil, fn.Assoc[string, int], as)

	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Expected %v, found %v", m, m2)
	}

	as = as.Array().Sort(fn.OrderTupleAsc[string, int])

	fn.SeqTest(t, as).Is(
		fn.TupleOf("one", 1),
		fn.TupleOf("three", 3),
		fn.TupleOf("two", 2))
}

func TestSeqAssocWhere(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}
	as := fn.AssocOf(m).Where(func(t fn.Tuple[string, int]) bool {
		return t.Key() == "one"
	})

	m2 := fn.Into(nil, fn.Assoc[string, int], as)

	if !reflect.DeepEqual(map[string]int{"one": 1}, m2) {
		t.Errorf("Expected %v, found %v", m, m2)
	}
}

func TestSeqAssocSkip(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	as := fn.AssocOf(m).Skip(0)
	m2 := fn.Into(nil, fn.Assoc[string, int], as)
	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Expected %v, found %v", m, m2)
	}

	as = fn.AssocOf(m).Skip(1)
	m2 = fn.Into(nil, fn.Assoc[string, int], as)
	if !(reflect.DeepEqual(map[string]int{"one": 1}, m2) || reflect.DeepEqual(map[string]int{"two": 2}, m2)) {
		t.Errorf("Expected {'one':1} or {'two': 2}, found %v", m2)
	}

	as = fn.AssocOf(m).Skip(123)
	m2 = fn.Into(nil, fn.Assoc[string, int], as)
	if len(m2) != 0 {
		t.Errorf("Expected empty map, found %v", m2)
	}
}

func TestSeqAssocTake(t *testing.T) {
	t.Skip("TODO")
}

func TestSeqAssocTakeWhile(t *testing.T) {
	t.Skip("TODO")
}
