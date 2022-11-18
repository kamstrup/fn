package fn

type Func0 func()

type Func0Err func() error

type Func1[T any] func(T)

type Func1Err[T any] func(T) error

type Func2[S, T any] func(S, T)

type Func2Err[S, T any] func(S, T) error

type FuncMap[S, T any] func(S) T

type FuncMapErr[S, T any] func(S) (T, error)

type FuncCollect[S, T any] func(T, S) T

type FuncCollectErr[S, T any] func(T, S) (T, error)

type FuncSource[T any] func() T

type FuncSourceErr[T any] func() (T, error)

type Predicate[T any] func(T) bool

type PredicateErr[T any] func(T) (bool, error)

func Sum[T Arithmetic](into, t T) T {
	return into + t
}

func Append[T any](into []T, t T) []T {
	return append(into, t)
}

func Collect[S comparable, T any](seq Seq[S], collector FuncCollect[S, T], into T) T {
	seq.ForEach(func(s S) {
		into = collector(into, s)
	})
	return into
}

func CollectErr[S comparable, T comparable](seq Seq[S], collector FuncCollectErr[S, T], into T) (T, error) {
	var err error
	seq.TakeWhile(func(s S) bool {
		into, err = collector(into, s)
		if err != nil {
			return false
		}
		return true
	})
	return into, err
}
