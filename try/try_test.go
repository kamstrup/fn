package fntry

import (
	"errors"
	"testing"

	fntesting "github.com/kamstrup/fn/testing"
)

var theError = errors.New("the error")
var atTheDiscoError = ErrPanic{V: "at the disco"}

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
	opt := Call(returnTheError)
	fntesting.OptOf(t, opt).IsError(theError)

	opt = CallRecover(panicAtTheDisco)
	fntesting.OptOf(t, opt).IsError(atTheDiscoError)

	opt = Apply(returnTheErrorIfEven, 28)
	fntesting.OptOf(t, opt).IsError(theError)

	opt = ApplyRecover(panicAtTheDiscoIfEven, 28)
	fntesting.OptOf(t, opt).IsError(atTheDiscoError)
}

func TestTryOk(t *testing.T) {
	opt := Call(func() (int, error) { return 27, nil })
	fntesting.OptOf(t, opt).Is(27)

	opt = CallRecover(func() (int, error) { return 27, nil })
	fntesting.OptOf(t, opt).Is(27)

	opt = Apply(returnTheErrorIfEven, 27)
	fntesting.OptOf(t, opt).Is(27)

	opt = ApplyRecover(panicAtTheDiscoIfEven, 27)
	fntesting.OptOf(t, opt).Is(27)
}
