package examples

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/slice"
)

func TestSliceMapping1(t *testing.T) {
	// Here we multiply all numbers in a slice with 2
	nums := []int{1, 2, 3}
	evens := slice.Mapping(nums, func(i int) int { return i * 2 })
	fmt.Println("a few even numbers:", evens)
}

func TestSliceMapping2(t *testing.T) {
	// Here we convert all numbers in a slice to strings base 2
	nums := []int{1, 2, 3}
	binaryStrings := slice.Mapping(nums, func(i int) string {
		return strconv.FormatInt(int64(i), 2)
	})
	fmt.Println("a few binary numbers as strings:", binaryStrings)
}

func TestSliceToMap(t *testing.T) {
	// Here we build a map that maps each word to their length
	words := []string{"one", "world"}
	strLens := slice.ToMap(words, func(word string) (string, int) {
		return word, len(word)
	})
	fmt.Println("word lengths:", strLens)
}

func TestSliceReduce(t *testing.T) {
	// Here join some strings via the fx.Reduce function
	// Note that the seq.FuncCollect functions work with fx.Reduce as well
	words := []string{"one", "world"}
	joinedWords := slice.Reduce(seq.MakeString, nil, words) // note: returns a strings.Builder
	fmt.Println("joined words:", joinedWords.String())
}

func TestSliceFromSeq(t *testing.T) {
	// Here we show how seq.Slice[T] can be used directly in the functions from the 'slice' package.
	shortWords := seq.SliceOfArgs("one", "world", "is", "not", "enough").
		Where(func(s string) bool { return len(s) <= 3 }).
		ToSlice()

	expected := []string{"one", "is", "not"}

	// shortWords is a seq.Slice[string], but we can still use it directly as a []string here
	asExpected := slice.Equal(expected, shortWords)

	fmt.Printf("are the short words %v as expected: %v\n", expected, asExpected)
}
