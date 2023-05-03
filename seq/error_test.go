package seq_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestError(t *testing.T) {
	theError := errors.New("hello")

	errSeq := seq.ErrorOf[any](theError)

	if err := errSeq.ForEach(func(a any) { t.Error("ForEach should not be called") }).Error(); err != theError {
		t.Error("ForEach must return theError")
	}

	if err := errSeq.ForEachIndex(func(i int, a any) { t.Error("ForEachIndex should not be called") }).Error(); err != theError {
		t.Error("ForEachIndex must return theError")
	}

	opt, cpy := errSeq.First()
	if cpy != errSeq {
		t.Errorf("First must return the errorSeq itself again")
	}

	_, errCpy := opt.Return()
	if errCpy != theError {
		t.Errorf("First must return an Opt that yield the original error")
	}
}

func TestErrorSuite(t *testing.T) {
	err := errors.New("hello")
	fntesting.SuiteOf(t, func() seq.Seq[int] { return seq.ErrorOf[int](err) }).IsEmpty()
}
