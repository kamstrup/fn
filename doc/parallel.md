Channels and Goroutines
====

seq.Chan
----
`seq.Chan[T]` is just a type wrapper for `chan T`, so any builtin function or syntax that applies to channels
can be used with `seq.Chan`. Eg

* `make(seq.Chan[string], 10)`
* `len(myChan)`, `cap(myChan)`
* `close(myChan)`
* for-range loops

Parallel Execution
----
You can execute a Seq in N goroutines mapping the results into a new Seq with `seq.Go()`:
```go
func fetchItem(id int) Opt[T] {
   // do something slow and calculate t
   return seq.OptOf(t) // or maybe an error
}

// Execute fetchItem of 1027 ids in 100 parallel goroutines
ids := seq.RangeOf(0, 1027)
result := seq.Go(ids, 100, fetchItem)

// result is a Seq[Opt[T]], let's print the successes and errors 
result.ForEach(func (opt Opt[T]) {
   t, err := opt.Error()
   if err != nil {
      fmt.Println("Oh no, an error!", err)
   } else {
      fmt.Println("Nice, got one T:", t)
   }
})
```
