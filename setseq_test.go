package fn

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	s := SetOf(map[string]struct{}{"foo": {}, "bar": {}})
	st := SeqTest(t, s)
	st.LenIs(2)

	exp := []string{"bar", "foo"}
	res := s.Array().Sort(OrderAsc[string]).AsSlice()
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("Expected %v, got %v", exp, res)
	}
}
