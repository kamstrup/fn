package opt

import (
	"sync"
)

// Future represents a result that will appear at some point in the future.
type Future[T any] struct {
	val Opt[T]
	wg  sync.WaitGroup
}

// Promise executes a function in a goroutine and returns a Future that can be used to wait for the result.
// The exec function *must* ensure that the resolve function is called exactly once.
// The promise does not have to be resolved by the exec function itself. It is allowed to pass the resolve function
// to other functions and goroutines.
//
// Calling resolve multiple times will cause a panic. Failing to call resolve will cause Future.Await to hang forever.
//
// If the exec function panics it will be recovered and Future.Await will return an option with an ErrPanic.
func Promise[T any](exec func(resolve func(Opt[T]))) *Future[T] {
	res := &Future[T]{}

	res.wg.Add(1)
	go func() {
		defer res.recoverPanic()
		exec(res.resolve)
	}()

	return res
}

func (fut *Future[T]) resolve(val Opt[T]) {
	fut.val = val
	fut.wg.Done()
}

func (fut *Future[T]) recoverPanic() {
	if e := recover(); e != nil {
		fut.val = ErrorOf[T](ErrPanic{
			V: e,
		})
	}
}

// Await blocks until the result is ready and returns it.
// If the result is already available the function returns immediately.
// It is valid to call from any goroutine and as many times as you like.
func (fut *Future[T]) Await() Opt[T] {
	fut.wg.Wait()
	return fut.val
}

// Then starts another promise result of this future is ready.
// The chained promise is started whether the first result is an error or not.
// If you need to change the type of the result you must use PromiseThen.
func (fut *Future[T]) Then(exec func(firstResult Opt[T], resolve func(Opt[T]))) *Future[T] {
	return PromiseThen[T, T](fut, exec)
}

// PromiseThen starts another promise when the result of a future is ready.
// The chained promise is started whether the first result is an error or not.
func PromiseThen[S, T any](first *Future[S], exec func(firstResult Opt[S], resolve func(Opt[T]))) *Future[T] {
	fut := &Future[T]{}

	fut.wg.Add(1)
	go func() {
		defer fut.recoverPanic()
		exec(first.Await(), fut.resolve)
	}()

	return fut
}
