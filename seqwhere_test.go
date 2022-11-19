package fn

import "testing"

func isOdd(i int) bool {
	return i%2 == 1
}

func isEven(i int) bool {
	return i%2 == 0
}

func TestWhere(t *testing.T) {
	odds := ArrayOfArgs(1, 2, 3, 4).Where(isOdd)
	SeqTest(t, odds).Is(1, 3)
}

func TestWhereSkip(t *testing.T) {
	odds := ArrayOfArgs(1, 2, 3, 4, 5, 6).Where(isOdd).Skip(2)
	SeqTest(t, odds).Is(5)

	odds = ArrayOfArgs(1, 2, 3, 4).Where(isOdd).Skip(0)
	SeqTest(t, odds).Is(1, 3)

	odds = ArrayOfArgs(1, 2, 3, 4).Where(isOdd).Skip(100)
	SeqTest(t, odds).IsEmpty()
}

func TestWhereTake(t *testing.T) {
	odds, _ := SourceOf(NumbersFrom(0)).Where(isOdd).Take(4)
	SeqTest(t, odds.Seq()).Is(1, 3, 5, 7)

	odds, _ = SourceOf(NumbersFrom(0)).Where(isOdd).Take(0)
	SeqTest(t, odds.Seq()).IsEmpty()

	odds, _ = SourceOf(NumbersFrom(0)).Where(isOdd).Take(1)
	SeqTest(t, odds.Seq()).Is(1)
}

func TestWhereTakeWhile(t *testing.T) {
	odds, _ := SourceOf(NumbersFrom(0)).Where(isOdd).TakeWhile(func(i int) bool { return i < 5 })
	SeqTest(t, odds.Seq()).Is(1, 3)

	odds, _ = SourceOf(NumbersFrom(0)).Where(isOdd).TakeWhile(func(_ int) bool { return false })
	SeqTest(t, odds.Seq()).IsEmpty()

	odds, _ = SourceOf(NumbersFrom(0)).Where(isOdd).TakeWhile(func(i int) bool { return i < 2 })
	SeqTest(t, odds.Seq()).Is(1)
}

func TestWhereWhile(t *testing.T) {
	odds := SourceOf(NumbersFrom(0)).Where(isOdd).While(func(i int) bool { return i < 5 })
	SeqTest(t, odds).Is(1, 3)

	odds = SourceOf(NumbersFrom(0)).Where(isOdd).While(func(_ int) bool { return false })
	SeqTest(t, odds).IsEmpty()

	odds = SourceOf(NumbersFrom(0)).Where(isOdd).While(func(i int) bool { return i < 2 })
	SeqTest(t, odds).Is(1)
}
