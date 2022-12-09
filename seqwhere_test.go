package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
)

func isOdd(i int) bool {
	return i%2 == 1
}

func isEven(i int) bool {
	return i%2 == 0
}

func TestWhere(t *testing.T) {
	odds := fn.ArrayOfArgs(1, 2, 3, 4).Where(isOdd)
	fn.SeqTest(t, odds).Is(1, 3)
}

func TestWhereSkip(t *testing.T) {
	odds := fn.ArrayOfArgs(1, 2, 3, 4, 5, 6).Where(isOdd).Skip(2)
	fn.SeqTest(t, odds).Is(5)

	odds = fn.ArrayOfArgs(1, 2, 3, 4).Where(isOdd).Skip(0)
	fn.SeqTest(t, odds).Is(1, 3)

	odds = fn.ArrayOfArgs(1, 2, 3, 4).Where(isOdd).Skip(100)
	fn.SeqTest(t, odds).IsEmpty()
}

func TestWhereTake(t *testing.T) {
	odds, _ := fn.NumbersFrom(0).Where(isOdd).Take(4)
	fn.SeqTest(t, odds.Seq()).Is(1, 3, 5, 7)

	odds, _ = fn.NumbersFrom(0).Where(isOdd).Take(0)
	fn.SeqTest(t, odds.Seq()).IsEmpty()

	odds, _ = fn.NumbersFrom(0).Where(isOdd).Take(1)
	fn.SeqTest(t, odds.Seq()).Is(1)
}

func TestWhereTakeWhile(t *testing.T) {
	odds, _ := fn.NumbersFrom(0).Where(isOdd).TakeWhile(func(i int) bool { return i < 5 })
	fn.SeqTest(t, odds.Seq()).Is(1, 3)

	odds, _ = fn.NumbersFrom(0).Where(isOdd).TakeWhile(func(_ int) bool { return false })
	fn.SeqTest(t, odds.Seq()).IsEmpty()

	odds, _ = fn.NumbersFrom(0).Where(isOdd).TakeWhile(func(i int) bool { return i < 2 })
	fn.SeqTest(t, odds.Seq()).Is(1)
}

func TestWhereWhile(t *testing.T) {
	odds := fn.NumbersFrom(0).Where(isOdd).While(func(i int) bool { return i < 5 })
	fn.SeqTest(t, odds).Is(1, 3)

	odds = fn.NumbersFrom(0).Where(isOdd).While(func(_ int) bool { return false })
	fn.SeqTest(t, odds).IsEmpty()

	odds = fn.NumbersFrom(0).Where(isOdd).While(func(i int) bool { return i < 2 })
	fn.SeqTest(t, odds).Is(1)
}
