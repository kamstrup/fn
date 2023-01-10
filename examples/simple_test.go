package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kamstrup/fn"
)

func TestExampleSimple(t *testing.T) {
	// This example we find words starting with "S", lowercase them and print them.
	words := fn.ArrayOfArgs("Reflection", "Collection", "Stream", "Sock").
		Where(func(s string) bool { return strings.HasPrefix(s, "S") }).
		Map(strings.ToLower).
		Map(func(s string) string {
			fmt.Println(s) // print as side effect of executing the 'words' seq
			return s
		})

	// Note: 'words' is still lazy so nothing has been done yet,
	// and the Println() statements follow after this Println()
	fmt.Println("Printing names starting with 'S'...")
	fn.Do(words)

	// If we wanted to execute immediately we could have replaced the last .Map() with .ForEach()
}

func TestExampleContains(t *testing.T) {
	// In this example we examine a sequence of names and checks if it contains "lisa"
	names := fn.ArrayOfArgs("John", "Bobby", "Lisa").
		Map(strings.ToLower)
	hasLisa := fn.Any(names, fn.Is("lisa"))
	fmt.Println(names, "lower-cased contains 'lisa':", hasLisa)
}

func TestExampleUserIndexes(t *testing.T) {
	// In this example we examine a sequence of usernames, and record the index of each occurrence
	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := fn.ZipOf[string, int](names, fn.NumbersFrom(0))
	userIndexes := fn.Into(nil, fn.GroupBy[string, int], tups)
	fmt.Println("Indexes of user names from", names, "\n", userIndexes)
}

func TestExampleUserSerial(t *testing.T) {
	// In this example we examine a sequence of usernames,
	// skip empty usernames and assign serial number to each unique user.
	names := fn.ArrayOfArgs("bob", "alan", "bob", "scotty", "", "bob", "alan", "").
		Where(fn.IsNonZero[string])
	tups := fn.ZipOf[string, int](names, fn.Constant(-1)) // the tuple seq is needed for UpdateAssoc
	serial := 0
	userSerials := fn.Into(nil, fn.UpdateAssoc[string, int](func(oldSerial, _ int) int {
		if oldSerial == 0 {
			serial += 1
			return serial
		}
		return oldSerial
	}), tups)
	fmt.Println("User serials from", names, "\n", userSerials)
}
