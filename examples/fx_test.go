package examples

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/fx"
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
	// Here join some strings via the fx.Into function
	words := []string{"one", "world"}
	joinedWords := fx.Into(nil, fn.MakeString, words) // note: returns a strings.Builder
	fmt.Println("joined words:", joinedWords.String())
}
