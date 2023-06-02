package seq_test

import (
	"testing"

	"github.com/kamstrup/fn/seq"
)

func TestTupleEquals(t *testing.T) {
	if !seq.TupleOf(1, 2).Equals(seq.TupleOf(1, 2)) {
		t.Fatal("tuples (1,2) and (1,2) must be equal")
	}

	if seq.TupleOf(1, 2).Equals(seq.TupleOf(1, 3)) {
		t.Fatal("tuples (1,3) must not be equal")
	}

	if seq.TupleOf(2, 2).Equals(seq.TupleOf(1, 2)) {
		t.Fatal("tuples (2,2) and (1, 2) must not be equal")
	}
}
