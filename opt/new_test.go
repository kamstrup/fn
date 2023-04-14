package opt

import (
	"errors"
	"reflect"
	"testing"
)

var theError = errors.New("the error")
var atTheDiscoError = ErrPanic{V: "at the disco"}

func isError[T any](t *testing.T, opt Opt[T], err error) {
	if opt.Ok() {
		t.Fatalf("expected error, opt was ok")
	} else {
		if opt.Error() != err {
			t.Fatalf("unexpected error: %v", opt.Error())
		}
	}
}

func is[T any](t *testing.T, opt Opt[T], val T) {
	if err := opt.Error(); err != nil {
		t.Fatalf("expacted valid opt, got: %v", err)
	} else {
		if !reflect.DeepEqual(val, opt.Must()) {
			t.Fatalf("bad value: %v", opt.Must())
		}
	}
}

func returnTheError() (int, error) {
	return 27, theError
}

func panicAtTheDisco() (int, error) {
	panic("at the disco")
}

func returnTheErrorIfEven(i int8) (int, error) {
	if i%2 == 0 {
		return int(i), theError
	}
	return int(i), nil
}

func panicAtTheDiscoIfEven(i int8) (int, error) {
	if i%2 == 0 {
		panic("at the disco")
	}
	return int(i), nil
}

func TestTryError(t *testing.T) {
	opt := Returning(returnTheError())
	isError(t, opt, theError)

	opt = Recovering(panicAtTheDisco)
	isError(t, opt, atTheDiscoError)

	opt = RecoveringMapper(panicAtTheDiscoIfEven)(28)
	isError(t, opt, atTheDiscoError)
}

func TestTryOk(t *testing.T) {
	opt := Returning(27, nil)
	is(t, opt, 27)

	opt = Recovering(func() (int, error) { return 27, nil })
	is(t, opt, 27)

	opt = RecoveringMapper(panicAtTheDiscoIfEven)(27)
	is(t, opt, 27)
}
