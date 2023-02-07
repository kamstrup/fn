package seq_test

import (
	"reflect"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestSet(t *testing.T) {
	s := seq.SetOf(map[string]struct{}{"foo": {}, "bar": {}})
	st := fntesting.TestOf(t, s)
	st.LenIs(2)

	exp := []string{"bar", "foo"}
	res := s.ToSlice().Sort(seq.OrderAsc[string])
	if !reflect.DeepEqual(exp, []string(res)) {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}
