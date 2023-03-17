package seq

import (
	"sync"

	"github.com/kamstrup/fn/opt"
)

// Go is like a parallelized eager version of MappingOf.
// It runs numJobs parallel goroutines of the task() function,
// and returns a Seq[T] that receives the results as they come in.
// This function returns immediately and the results are computed in the background.
//
// If you just want to execute a Seq for the sake of triggering side effects
// you can use MapOf combined with the Do function.
func Go[S, T any](seq Seq[S], numJobs int, task FuncMap[S, T]) Seq[T] {
	chS := make(chan S, numJobs)
	chT := make(chan T, numJobs)
	wg := sync.WaitGroup{}
	wg.Add(numJobs)

	// Start N goroutines doing work off chS, converting S -> T
	for i := 0; i < numJobs; i++ {
		go func() {
			defer wg.Done()
			for s := range chS {
				chT <- task(s)
			}
		}()
	}

	// We need to capture the tail of seq.ForEach to detect errors,
	// which is tricky because we only have the tail *after* the input chan is drained.
	tailPromise := opt.Promise(func(resolve func(opt.Opt[T])) {
		// Start pumping work into chS and indicate completion with close()
		seqTail := seq.ForEach(func(s S) {
			chS <- s
		})

		close(chS)

		// Wait until all tasks are done processing,
		// then try to send the error if there is one,
		// and finally close chT to mark the end of the output Seq[T]
		wg.Wait()
		close(chT)

		if err := Error(seqTail); err != nil {
			resolve(opt.ErrorOf[T](err))
		} else {
			resolve(opt.Empty[T]())
		}
	})

	// We create a tail that lazily resolves the tailPromise.
	// The tail will capture any error returned by seq.ForEach.
	tail := ValuesOf(SourceOf(tailPromise.Await).Limit(1))
	return ConcatOf(ChanOf(chT), tail)
}
