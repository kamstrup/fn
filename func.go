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

func Count[T Arithmetic](into, t T) T {
	return into + 1
}

func Append[T any](into []T, t T) []T {
	return append(into, t)
}

func Assoc[K, V comparable](into map[K]V, t Tuple[K, V]) map[K]V {
	into[t.Key()] = t.Value()
	return into
}

func Collect[S any, T any](seq Seq[S], collector FuncCollect[S, T], into T) T {
	seq.ForEach(func(s S) {
		into = collector(into, s)
	})
	return into
}

func CollectErr[S any, T any](seq Seq[S], collector FuncCollectErr[S, T], into T) (T, error) {
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
