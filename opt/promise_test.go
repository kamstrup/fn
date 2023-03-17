package opt

import (
	"math/rand"
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	futures := make([]*Future[int], 100)
	for i := 0; i < len(futures); i++ {
		j := i
		futures[i] = Promise(func(resolve func(opt Opt[int])) {
			time.Sleep(time.Duration(rand.Int63n(1000)))
			if j%10 == 0 {
				resolve(ErrorOf[int](theError))
			} else {
				resolve(Of(j))
			}
		})
	}

	for i, f := range futures {
		val, err := f.Await().Return()
		if i%10 == 0 {
			if err != theError {
				t.Fatalf("expected %q, found: %q", theError, err)
			}
		} else {
			if i != val {
				t.Fatalf("bad value, expected %d, was: %d", i, val)
			}
		}
	}
}

func TestPromiseThen(t *testing.T) {
	futures := make([]*Future[int], 100)
	for i := 0; i < len(futures); i++ {
		j := i
		futures[i] = Promise(func(resolve func(opt Opt[int])) {
			time.Sleep(time.Duration(rand.Int63n(1000)))
			if j%10 == 0 {
				resolve(ErrorOf[int](theError))
			} else {
				resolve(Of(j))
			}
		}).Then(func(numOpt Opt[int], resolve func(Opt[int])) {
			if num, err := numOpt.Return(); err != nil {
				resolve(numOpt)
			} else {
				resolve(Of(num * 2))
			}
		})
	}

	for i, f := range futures {
		val, err := f.Await().Return()
		if i%10 == 0 {
			if err != theError {
				t.Fatalf("expected %q, found: %q", theError, err)
			}
		} else {
			if i*2 != val {
				t.Fatalf("bad value, expected %d, was: %d", i*2, val)
			}
		}
	}
}
