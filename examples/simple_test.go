package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kamstrup/fn"
)

func TestExampleSimple(t *testing.T) {
	names := fn.ArrayOfArgs("Reflection", "Collection", "Stream", "Sock").
		Where(func(s string) bool { return strings.HasPrefix(s, "S") }).
		Map(strings.ToLower).
		Array()
	fmt.Println(names)
}

func TestExampleContains(t *testing.T) {
	names := fn.ArrayOfArgs("John", "Bobby", "Lisa").
		Map(strings.ToLower)
	hasLisa := fn.Any(names, fn.Is("lisa"))
	fmt.Println(names, "lower-cased contains 'lisa':", hasLisa)
}
