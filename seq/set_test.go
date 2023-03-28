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

func TestSetContains(t *testing.T) {
	s := seq.SetAsArgs("foo", "bar")
	if sz, ok := s.Len(); !ok || sz != 2 {
		t.Errorf("unexpected set len: %d", sz)
	}

	seq.SliceOfArgs("foo", "bar").ForEach(func(str string) {
		if !s.Contains(str) {
			t.Errorf("set must contain %s", str)
		}
	})

	seq.SliceOfArgs("hello", "world").ForEach(func(str string) {
		if s.Contains(str) {
			t.Errorf("set must not contain %s", str)
		}
	})
}

func TestSetUnion(t *testing.T) {
	expectElems := seq.SliceOfArgs("foo", "bar", "boo")
	s1 := seq.SetAsArgs("foo", "bar")
	s2 := seq.SetAsArgs("boo", "bar")

	u1 := s1.Union(s2)
	u2 := s2.Union(s1)

	seq.SliceOfArgs(u1, u2).ForEach(func(u seq.Seq[string]) {
		su := seq.Reduce(seq.MakeSet[string], nil, u).Or(nil)
		if len(su) != 3 {
			t.Errorf("set length must be 3: %v", su.ToSlice())
		}
		if !seq.All(expectElems, su.Contains) {
			t.Errorf("set must contain all elements: %v, found %v", expectElems, su.ToSlice())
		}
	})

	selfU := s1.Union(s1) // self union
	selfSlice := selfU.ToSlice().Sort(seq.OrderAsc[string])
	if !reflect.DeepEqual(seq.SliceAsArgs("bar", "foo"), selfSlice) {
		t.Errorf("set must be [foo, bar], found: %v", selfSlice)
	}
}

func TestSetIntersect(t *testing.T) {
	expectElems := seq.SliceOfArgs("bar")
	s1 := seq.SetAsArgs("foo", "bar")
	s2 := seq.SetAsArgs("boo", "bar")

	i1 := s1.Intersect(s2)
	i2 := s2.Intersect(s1)

	seq.SliceOfArgs(i1, i2).ForEach(func(i seq.Seq[string]) {
		si := seq.Reduce(seq.MakeSet[string], nil, i).Or(nil)
		if len(si) != 1 {
			t.Errorf("set length must be 3: %v", si.ToSlice())
		}
		if !seq.All(expectElems, si.Contains) {
			t.Errorf("set must contain all elements: %v, found %v", expectElems, si.ToSlice())
		}
	})

	selfI := s1.Intersect(s1) // self intersect
	selfSlice := selfI.ToSlice().Sort(seq.OrderAsc[string])
	if !reflect.DeepEqual(seq.SliceAsArgs("bar", "foo"), selfSlice) {
		t.Errorf("set must be [bar, foo], found: %v", selfSlice)
	}
}
