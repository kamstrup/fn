package fn

import "sync"

// Go is like a parallelized version, eager, version of MapOf.
// It runs numJobs parallel goroutines of the task() function,
// and returns a Seq[T] that receives the results as they come in.
// This function returns immediately and the results are computed in the background.
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

	// The controlling goroutine
	go func() {
		// Start pumping work into chS and indicate completion with close()
		seq.ForEach(func(s S) {
			chS <- s
		})
		close(chS)

		// Wait until all tasks are done processing,
		// then close cht to mark the end of the output Seq[T]
		wg.Wait()
		close(chT)
	}()

	return ChanOf(chT)
}
