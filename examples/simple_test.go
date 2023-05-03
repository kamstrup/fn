package examples

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

func TestExampleSimple(t *testing.T) {
	// This example we find words starting with "S", lowercase them and print them.
	words := seq.SliceOfArgs("Reflection", "Collection", "Stream", "Sock").
		Where(func(s string) bool { return strings.HasPrefix(s, "S") }).
		Map(strings.ToLower).
		Map(func(s string) string {
			fmt.Println(s) // print as side effect of executing the 'words' seq
			return s
		})

	// Note: 'words' is still lazy so nothing has been done yet,
	// and the Println() statements follow after this Println()
	fmt.Println("Printing names starting with 'S'...")
	seq.Do(words)
	// Prints:
	// Printing names starting with 'S'...
	// stream
	// sock

	// If we wanted to execute immediately we could have replaced the last .Map() with .ForEach()
}

func TestExampleContains(t *testing.T) {
	// In this example we examine a sequence of names and checks if it contains "lisa"
	names := seq.SliceOfArgs("John", "Bobby", "Lisa").
		Map(strings.ToLower)
	hasLisa := seq.Any(names, seq.Is("lisa"))
	fmt.Println(names, "lower-cased contains 'lisa':", hasLisa) // true
}

func TestExampleUserIndexes(t *testing.T) {
	// In this example we examine a sequence of usernames, and record the index of each occurrence
	names := seq.SliceOfArgs("bob", "alan", "bob", "scotty", "bob", "alan")
	tups := seq.ZipOf[string, int](names, seq.RangeFrom(0))
	userIndexes := seq.Reduce(seq.GroupBy[string, int], nil, tups).Or(nil)
	fmt.Println("Indexes of user names from", names, "\n", userIndexes)
	// Prints:
	// Indexes of user names from [bob alan bob scotty bob alan]
	// map[alan:[1 5] bob:[0 2 4] scotty:[3]]
}

func TestExampleUserSerial(t *testing.T) {
	// In this example we examine a sequence of usernames,
	// skip empty usernames and assign serial number to each unique user.
	names := seq.SliceOfArgs("bob", "alan", "bob", "scotty", "", "bob", "alan", "").
		Where(seq.IsNonZero[string])
	tups := seq.ZipOf[string, int](names, seq.Constant(-1)) // the tuple seq is needed for UpdateMap
	serial := 0
	userSerials := seq.Reduce(seq.UpdateMap[string, int](func(oldSerial, _ int) int {
		if oldSerial == 0 {
			serial += 1
			return serial
		}
		return oldSerial
	}), nil, tups).
		Or(nil)

	fmt.Println("User serials from", names, "\n", userSerials)
	// Prints:
	// User serials from {[bob alan bob scotty  bob alan ] 0x4fcd20}
	// map[alan:2 bob:1 scotty:3]
}

func TestExampleParallelDownloader(t *testing.T) {
	// In this example we simulate a massive parallel download
	// of a bunch of files named numerically 0.txt ... 1027.txt.
	// Keeping 100 items in-flight all the time
	fetchItem := func(num int) int {
		fmt.Printf("Downloading %d.txt\n", num) // not really, just a dummy
		return num
	}

	ids := seq.RangeOf(0, 1027)
	tasks := seq.Go(ids, 100, fetchItem)
	result := seq.Do(tasks)

	if err := result.Error(); err != nil {
		t.Fatal(err) // Not going to happen in this test, but might in real apps
	}
	fmt.Println("All done")
}

func TestExampleMapType(t *testing.T) {
	// In this example we show how a seq.Map can be created as a literal an used
	// as a normal Go map (same is true for ChanAs, SliceAs, SetAs etc)
	// and also as a Seq
	myMap := seq.Map[string, int]{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
	}

	fmt.Println("Number of elements in myMap:", len(myMap)) // prints 2

	for k, v := range myMap {
		fmt.Println("Ranging over myMap:", k, v)
	}

	// We can also use myMap as a Seq
	// Let's find entries with even values and collect the keys for them:
	evenEntries := myMap.Where(func(kv seq.Tuple[string, int]) bool {
		return kv.Value()%2 == 0
	})
	evenKeys := seq.MappingOf(evenEntries, seq.TupleKey[string, int]).ToSlice()
	fmt.Println("Map keys with even values:", evenKeys)

	// Prints:
	// Number of elements in myMap: 4
	// Ranging over myMap: one 1
	// Ranging over myMap: two 2
	// Ranging over myMap: three 3
	// Ranging over myMap: four 4
	// Map keys with even values: [two four]
}

func TestExampleMapWithErrors(t *testing.T) {
	// In this example we use a mapping function that can also return an error.
	// We try to parse a series of strings as integers and if it fails we silently drop them
	strInts := seq.SliceOfArgs("1", "two", "3")
	ints := seq.MappingOf(strInts, opt.Mapper(strconv.Atoi)).
		Where(opt.Ok[int]).
		ToSlice()

	fmt.Println("These ints parsed ok:", ints)
}

func TestExampleSourceWithErrors(t *testing.T) {
	// In this example we use a mapping function that can also return an error.
	// We try to parse a series of strings as integers and if it fails we silently drop them
	numTooBigError := errors.New("num too big!")
	i := 0
	ints, tail := seq.SourceOf(opt.Caller(func() (int, error) {
		i++
		if i > 3 {
			return 0, numTooBigError
		}
		return i, nil
	})).TakeWhile(opt.Ok[int])

	fmt.Println("These ints are ok:", ints)
	fmt.Println("Tail ended up as:", tail)
}

func TestExampleTeams(t *testing.T) {
	blueTeam := seq.SliceOfArgs("Alan", "Betty")
	redTeam := seq.SliceOfArgs("Maria", "Bob")

	// Let's create a set of names for the people on the blue and red teams
	allTeamMembers := seq.ConcatOf(blueTeam, redTeam)
	nameSet := seq.Reduce(seq.MakeSet[string], nil, allTeamMembers).Or(nil)

	// We need 2 members for the green team, that are not already on the blue or red team
	greenTeam := seq.SliceOfArgs("Betty", "Maurice", "Bob", "Charles", "Inga").
		Where(seq.Not(nameSet.Contains)).
		Limit(2).
		ToSlice()

	if len(greenTeam) != 2 {
		panic("not enough team members for the green team")
	}

	// Members on the green team are assigned player numbers starting from 10
	greenTeamNumbers := seq.ZipOf[int, string](seq.RangeFrom(10), greenTeam)
	greenTeamNumbers.ForEach(func(member seq.Tuple[int, string]) {
		fmt.Println("Name:", member.Value(), "Number:", member.Key())
	})
	// Prints:
	// Name: Maurice Number: 10
	// Name: Charles Number: 11
}
