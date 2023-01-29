package fn_test

import (
	"errors"
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func isOdd(i int) bool {
	return i%2 == 1
}

func isEven(i int) bool {
	return i%2 == 0
}

func TestWhere(t *testing.T) {
	odds := fn.ArrayOfArgs(1, 2, 3, 4).Where(isOdd)
	fntesting.TestOf(t, odds).Is(1, 3)
}

func TestWhereSkip(t *testing.T) {
	odds := fn.ArrayOfArgs(1, 2, 3, 4, 5, 6).Where(isOdd).Skip(2)
	fntesting.TestOf(t, odds).Is(5)

	odds = fn.ArrayOfArgs(1, 2, 3, 4).Where(isOdd).Skip(0)
	fntesting.TestOf(t, odds).Is(1, 3)

	odds = fn.ArrayOfArgs(1, 2, 3, 4).Where(isOdd).Skip(100)
	fntesting.TestOf(t, odds).IsEmpty()
}

func TestWhereTake(t *testing.T) {
	odds, _ := fn.RangeFrom(0).Where(isOdd).Take(4)
	fntesting.TestOf(t, odds.Seq()).Is(1, 3, 5, 7)

	odds, _ = fn.RangeFrom(0).Where(isOdd).Take(0)
	fntesting.TestOf(t, odds.Seq()).IsEmpty()

	odds, _ = fn.RangeFrom(0).Where(isOdd).Take(1)
	fntesting.TestOf(t, odds.Seq()).Is(1)
}

func TestWhereTakeWhile(t *testing.T) {
	odds, _ := fn.RangeFrom(0).Where(isOdd).TakeWhile(func(i int) bool { return i < 5 })
	fntesting.TestOf(t, odds.Seq()).Is(1, 3)

	odds, _ = fn.RangeFrom(0).Where(isOdd).TakeWhile(func(_ int) bool { return false })
	fntesting.TestOf(t, odds.Seq()).IsEmpty()

	odds, _ = fn.RangeFrom(0).Where(isOdd).TakeWhile(func(i int) bool { return i < 2 })
	fntesting.TestOf(t, odds.Seq()).Is(1)
}

func TestWhereWhile(t *testing.T) {
	odds := fn.RangeFrom(0).Where(isOdd).While(func(i int) bool { return i < 5 })
	fntesting.TestOf(t, odds).Is(1, 3)

	odds = fn.RangeFrom(0).Where(isOdd).While(func(_ int) bool { return false })
	fntesting.TestOf(t, odds).IsEmpty()

	odds = fn.RangeFrom(0).Where(isOdd).While(func(i int) bool { return i < 2 })
	fntesting.TestOf(t, odds).Is(1)
}

func TestWhereError(t *testing.T) {
	theError := errors.New("the error")
	wh := fn.ErrorOf[int](theError).Where(fn.IsNonZero[int])

	if err := fn.Error(wh); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
