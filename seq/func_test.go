package seq_test

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	fnmath "github.com/kamstrup/fn/math"
	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestZeroes(t *testing.T) {
	fntesting.TestOf(t, seq.MappingOf(seq.SliceOfArgs(-1, 0, 1, 10), seq.IsZero[int])).Is(false, true, false, false)
	fntesting.TestOf(t, seq.MappingOf(seq.SliceOfArgs(-1, 0, 1, 10), seq.IsNonZero[int])).Is(true, false, true, true)
}

func TestCollectCount(t *testing.T) {
	arr := seq.SliceOfArgs[int](1, 2, 3)

	count := seq.Reduce(seq.Count[int], 0, arr)
	if count.Must() != 3 {
		t.Errorf("expected count 3: %d", count)
	}

	count = seq.Reduce(seq.Count[int], 0, seq.Empty[int]())
	if !count.Empty() || count.Ok() {
		t.Errorf("expected empty count: %v", count)
	}
	if val, err := count.Return(); val == 27 || err != opt.ErrEmpty {
		t.Errorf("expected empty count: %v", count)
	}
}

func TestCollectAppend(t *testing.T) {
	arr := seq.SliceOfArgs[int](1, 2, 3)
	cpy := seq.Reduce(seq.Append[int], nil, arr)
	exp := []int{1, 2, 3}
	if !reflect.DeepEqual(cpy.Must(), exp) {
		t.Errorf("expected %v, got %v", exp, cpy)
	}

	ints := seq.Reduce(seq.Append[int], []int{27}, seq.Empty[int]())
	if !ints.Empty() || ints.Ok() {
		t.Errorf("expected empty min: %v", ints)
	}
	if val, err := ints.Return(); len(val) != 0 || err != opt.ErrEmpty {
		t.Errorf("expected empty ints: %v", ints)
	}
}

func TestCollectAssoc(t *testing.T) {
	oddNums := seq.SliceOfArgs(1, 2, 3).
		Where(func(i int) bool { return i%2 == 1 })

	arr := seq.MappingOf(oddNums, seq.TupleWithKey(func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	}))
	res := seq.Reduce(seq.MakeMap[string, int], nil, arr)
	exp := map[string]int{
		"1": 1, "3": 3,
	}
	if !reflect.DeepEqual(res.Must(), exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectSet(t *testing.T) {
	nums := seq.SliceOfArgs(1, 2, 2, 3, 1)
	res := seq.Reduce(seq.MakeSet[int], nil, nums)
	exp := map[int]struct{}{
		1: {}, 2: {}, 3: {},
	}
	if !reflect.DeepEqual(res.Must(), exp) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectString(t *testing.T) {
	strs := seq.SliceOfArgs("one", "two")
	res := seq.Reduce(seq.MakeString, nil, strs)
	exp := "onetwo"
	if exp != res.Must().String() {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectGroupBy(t *testing.T) {
	names := seq.SliceOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := seq.ZipOf[string, int](names, seq.RangeFrom(0))
	res := seq.Reduce(seq.GroupBy[string, int], nil, tups)
	exp := map[string][]int{
		"bob":    {0, 2, 4},
		"alan":   {1, 5},
		"scotty": {3},
	}

	if !reflect.DeepEqual(exp, res.Must()) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectUpdateAssoc(t *testing.T) {
	names := seq.SliceOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := seq.ZipOf[string, int](names, seq.Constant(1))
	res := seq.Reduce(seq.UpdateMap[string, int](fnmath.Sum[int]), nil, tups)
	exp := map[string]int{
		"bob":    3,
		"alan":   2,
		"scotty": 1,
	}

	if !reflect.DeepEqual(exp, res.Must()) {
		t.Errorf("expected %v, got %v", exp, res)
	}
}

func TestCollectUpdateArray(t *testing.T) {

	hellos := seq.SliceOfArgs(
		seq.TupleOf(1, "hello"), seq.TupleOf(2, "hej"),
		seq.TupleOf(1, "world"), seq.TupleOf(2, "verden"))

	res := seq.Reduce(seq.UpdateSlice[int, string](func(old, new_ string) string {
		return strings.TrimSpace(old + " " + new_)
	}), nil, hellos)

	exp := []string{
		"",
		"hello world",
		"hej verden",
	}

	if !reflect.DeepEqual(exp, res.Must()) {
		t.Errorf("expected %v, got %v", exp, res.Must())
	}
}

func TestCollectChan(t *testing.T) {
	orig := seq.SliceAsArgs(1, 2, 3)
	ch := seq.Reduce(seq.MakeChan[int], make(chan int, len(orig)), orig.Seq()).Must()
	close(ch)
	chSlice := seq.ChanOf(ch).ToSlice()

	if !reflect.DeepEqual(chSlice, orig) {
		t.Fatalf("bad result: %v", chSlice)
	}
}

func TestCollectError(t *testing.T) {
	theError := errors.New("the error")
	errSeq := seq.ErrorOf[int](theError)
	res := seq.Reduce(fnmath.Sum[int], 0, errSeq)

	if res.Error() != theError || res.Ok() {
		t.Errorf("expected 'the error': %v", res)
	}
}

func TestPredicates(t *testing.T) {
	if seq.LessThanZero(1) {
		t.Errorf("1 is not < 0")
	}
	if !seq.LessThanZero(-1) {
		t.Errorf("-1 is < 0")
	}

	if !seq.GreaterThanZero(1) {
		t.Errorf("1 is > 0")
	}
	if seq.GreaterThanZero(-1) {
		t.Errorf("-1 is not > 0")
	}

	if seq.Is("hello")("hej") {
		t.Errorf("hello != hej")
	}
	if !seq.IsNot("hello")("hej") {
		t.Errorf("hello != hej")
	}

	if seq.IsZero("hello") {
		t.Errorf("hello is non-zero")
	}
	if !seq.IsZero("") {
		t.Errorf("\"\" should be zero")
	}
	if !seq.IsNonZero("hello") {
		t.Errorf("hello is non-zero")
	}
	if seq.IsNonZero("") {
		t.Errorf("\"\" should be zero")
	}
}

func TestAny(t *testing.T) {
	if seq.Any(seq.SliceOfArgs(1, 2, 3), seq.IsZero[int]) {
		t.Fatal("should not find zero")
	}

	if !seq.Any(seq.SliceOfArgs(0, 1), seq.IsNonZero[int]) {
		t.Fatal("should find non-zero")
	}
}

func TestAll(t *testing.T) {
	if seq.All(seq.SliceOfArgs(0, 0, 1), seq.IsZero[int]) {
		t.Fatal("should find no.zero")
	}

	if !seq.All(seq.SliceOfArgs(0, 0), seq.IsZero[int]) {
		t.Fatal("should be all zeroes")
	}
}

func TestLast(t *testing.T) {
	o := seq.Last(seq.SliceOfArgs(0, 0, 1))
	if o.Error() != nil {
		t.Fatal("should not error", o.Error())
	}
	if o.Must() != 1 {
		t.Fatal("should be 1", o.Must())
	}

	o = seq.Last(seq.SliceOfArgs(0))
	if o.Error() != nil {
		t.Fatal("should not error", o.Error())
	}
	if o.Must() != 0 {
		t.Fatal("should be 0", o.Must())
	}

	o = seq.Last(seq.Empty[int]())
	if !o.Empty() {
		t.Fatal("should be empty", o.Must())
	}

	theError := errors.New("the error")
	o = seq.Last(seq.ErrorOf[int](theError))
	if !o.Empty() || o.Ok() {
		t.Fatal("should be empty", o.Must())
	}
	if o.Error() != theError {
		t.Fatal("should be 'the error'", o.Error())
	}
}
