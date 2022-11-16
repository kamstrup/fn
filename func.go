package fn

type Func0 func()

type Func0Err func() error

type Func1[T any] func(T)

type Func1Err[T any] func(T) error

type Func2[S, T any] func(S, T)

type Func2Err[S, T any] func(S, T) error

type FuncMap[S, T any] func(S) T

type FuncMapErr[S, T any] func(S) (T, error)

type FuncSource[T any] func() T

type FuncSourceErr[T any] func() (T, error)

type Predicate[T any] func(T) bool

type PredicateErr[T any] func(T) (bool, error)
