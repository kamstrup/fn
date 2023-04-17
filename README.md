# Fn(), Functional Programming for Golang

[![Go Report Card](https://goreportcard.com/badge/github.com/kamstrup/fn)](https://goreportcard.com/report/github.com/kamstrup/fn) [![PkgGoDev](https://pkg.go.dev/badge/github.com/kamstrup/fn)](https://pkg.go.dev/github.com/kamstrup/fn)

Fn is library for golang that enable you to blend functional programming techniques
with standard idiomatic Go code.

We are inspired by [Clojure](https://github.com/clojure/clojure),
[Vavr](https://github.com/vavr-io/vavr), and the Java Streams APIs that were
introduced back in Java 8, and want to provide something of similar spirit
that makes it even more fun to write Go code.

## Documentation
You will find [comprehensive documentation in our docs folder](doc), or you might want to dive
directly into our [simple examples](https://github.com/kamstrup/fn/blob/main/examples/simple_test.go).

## Examples
```go
import (
    "fmt"
	
    "github.com/kamstrup/fn/seq"
)

func printGreenTeam() {
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
    greenTeamNumbers.ForEach(func (member seq.Tuple[int, string]) {
        fmt.Println("Name:", member.Value(), "Number:", member.Key())
    })
    // Prints:
    // Name: Maurice Number: 10
    // Name: Charles Number: 11
}
```

## Performance
If the foundational functional data structures and algorithms is not done carefully,
execution speed and memory usage will suffer. Fn() is designed to make the best of what
the Go runtime can provide. Initial benchmarks puts it as a top contender among Golang
functional libraries. See benchmarks here https://github.com/mariomac/go-stream-benchmarks/pull/1
