package fn

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	err := errors.New("hello")

	errSeq := ErrorOf[any](err)
	if Error(errSeq) != err {
		t.Errorf("Error() on errorSeq must return the wrapped error")
	}

	errSeq.ForEach(func(a any) { t.Error("ForEach should not be called") })
	errSeq.ForEachIndex(func(i int, a any) { t.Error("ForEachIndex should not be called") })

	opt, cpy := errSeq.First()
	if cpy != errSeq {
		t.Errorf("First must return the errorSeq itself again")
	}

	_, errCpy := opt.Return()
	if errCpy != err {
		t.Errorf("First must return an Opt that yield the original error")
	}
}

func TestErrorSuite(t *testing.T) {
	err := errors.New("hello")
	SeqTestSuite(t, func() Seq[int] { return ErrorOf[int](err) }).IsEmpty()
}
