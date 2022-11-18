package fn

import (
	"reflect"
	"testing"
)

func TestSeqAssoc(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}
	as := AssocOf(m)

	SeqTest(t, as).LenIs(2)
	m2 := Into(nil, Assoc[string, int], as)

	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Expected %v, found %v", m, m2)
	}
}

func TestSeqAssocWhere(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}
	as := AssocOf(m).Where(func(t Tuple[string, int]) bool {
		return t.Key() == "one"
	})

	m2 := Into(nil, Assoc[string, int], as)

	if !reflect.DeepEqual(map[string]int{"one": 1}, m2) {
		t.Errorf("Expected %v, found %v", m, m2)
	}
}

func TestSeqAssocSkip(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	as := AssocOf(m).Skip(0)
	m2 := Into(nil, Assoc[string, int], as)
	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Expected %v, found %v", m, m2)
	}

	as = AssocOf(m).Skip(1)
	m2 = Into(nil, Assoc[string, int], as)
	if !(reflect.DeepEqual(map[string]int{"one": 1}, m2) || reflect.DeepEqual(map[string]int{"two": 2}, m2)) {
		t.Errorf("Expected {'one':1} or {'two': 2}, found %v", m2)
	}

	as = AssocOf(m).Skip(123)
	m2 = Into(nil, Assoc[string, int], as)
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
