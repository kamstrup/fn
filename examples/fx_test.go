package examples

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kamstrup/fn/fx"
	"github.com/kamstrup/fn/seq"
)

func TestFxMapSlice1(t *testing.T) {
	// Here we multiply all numbers in a slice with 2
	nums := []int{1, 2, 3}
	evens := fx.MapSlice(nums, func(i int) int { return i * 2 })
	fmt.Println("a few even numbers:", evens)
}

func TestFxMapSlice2(t *testing.T) {
	// Here we convert all numbers in a slice to strings base 2
	nums := []int{1, 2, 3}
	binaryStrings := fx.MapSlice(nums, func(i int) string {
		return strconv.FormatInt(int64(i), 2)
	})
	fmt.Println("a few binary numbers as strings:", binaryStrings)
}

func TestFxAssocSlice(t *testing.T) {
	// Here we build a map that maps each word to their length
	words := []string{"one", "world"}
	strLens := fx.AssocSlice(words, func(word string) (string, int) {
		return word, len(word)
	})
	fmt.Println("word lengths:", strLens)
}

func TestFxIntoStrings(t *testing.T) {
	// Here join some strings via the fx.Reduce function
	// Note that the seq.FuncCollect functions work with fx.Reduce as well
	words := []string{"one", "world"}
	joinedWords := fx.Reduce(seq.MakeString, nil, words) // note: returns a strings.Builder
	fmt.Println("joined words:", joinedWords.String())
}
