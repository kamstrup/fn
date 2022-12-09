package fn_test

import (
	"reflect"
	"testing"

	"github.com/kamstrup/fn"
)

func TestSet(t *testing.T) {
	s := fn.SetOf(map[string]struct{}{"foo": {}, "bar": {}})
	st := fn.SeqTest(t, s)
	st.LenIs(2)

	exp := []string{"bar", "foo"}
	res := s.Array().Sort(fn.OrderAsc[string]).AsSlice()
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}
