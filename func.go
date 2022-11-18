package fn

type Func0 func()

type Func0Err func() error

type Func1[T any] func(T)

type Func1Err[T any] func(T) error

type Func2[S, T any] func(S, T)

type Func2Err[S, T any] func(S, T) error

type FuncMap[S, T any] func(S) T

type FuncMapErr[S, T any] func(S) (T, error)

// FuncCollect is used to aggregate data and return the updated aggregation.
// Think of it as a generic form of the standard append() function in go.
// The order of the method arguments also follow the convention from append() and copy(),
// having the target destination as first argument.
type FuncCollect[S, T any] func(T, S) T

// FuncCollectErr is used to aggregate data and return the updated aggregation.
// Think of it as a generic form of the standard append() function in go.
// This is a variation of FuncCollect that can return an error,
// which should cause aggregation to stop.
type FuncCollectErr[S, T any] func(T, S) (T, error)

type FuncSource[T any] func() T

type FuncSourceErr[T any] func() (T, error)

type FuncLess[T any] func(T, T) bool

type Predicate[T any] func(T) bool

type PredicateErr[T any] func(T) (bool, error)

func Sum[T Arithmetic](into, t T) T {
	return into + t
}

func Count[T Arithmetic](into, _ T) T {
	return into + 1
}

func Append[T any](into []T, t T) []T {
	return append(into, t)
}

// Assoc is a FuncCollect that can take a Seq of Tuple elements and store them in a standard Go map.
func Assoc[K comparable, V any](into map[K]V, t Tuple[K, V]) map[K]V {
	if into == nil {
		into = make(map[K]V)
	}
	into[t.Key()] = t.Value()
	return into
}

// Set is a FuncCollect that can take a Seq of comparable values and store them in a standard Go set (map[]struct{}).
func Set[K comparable](into map[K]struct{}, k K) map[K]struct{} {
	if into == nil {
		into = make(map[K]struct{})
	}
	into[k] = struct{}{}
	return into
}

// OrderAsc is a FuncLess that can be used with Array.Sort
func OrderAsc[T Ordered](t1, t2 T) bool {
	return t1 < t2
}

// OrderDesc is a FuncLess that can be used with Array.Sort
func OrderDesc[T Ordered](t1, t2 T) bool {
	return t1 > t2
}

// OrderTupleAsc is a FuncLess that can be used with Array.Sort
func OrderTupleAsc[K Ordered, V any](t1, t2 Tuple[K, V]) bool {
	return t1.Key() < t2.Key()
}

// OrderTupleDesc is a FuncLess that can be used with Array.Sort
func OrderTupleDesc[K Ordered, V any](t1, t2 Tuple[K, V]) bool {
	return t1.Key() > t2.Key()
}

// NumbersFrom returns a FuncSource that starts from n and count one up on every invocation.
// You can for example use it with SourceOf()
func NumbersFrom(n int) FuncSource[int] {
	counter := n - 1
	return func() int {
		counter += 1
		return counter
	}
}

// Constant returns a FuncSource that produces the same value on every invocation.
// You can for example use it with SourceOf()
func Constant[T any](t T) FuncSource[T] {
	return func() T {
		return t
	}
}

// Into executes a Seq, collecting the results via a collection function.
// The method signature follows append() and copy() conventions,
// having the destination to put data into first.
func Into[S any, T any](into T, collector FuncCollect[S, T], seq Seq[S]) T {
	seq.ForEach(func(s S) {
		into = collector(into, s)
	})
	return into
}

// IntoErr is like Into() but the collector function can return an error,
// causing collection to stop and the error returned.
func IntoErr[S any, T any](into T, collector FuncCollectErr[S, T], seq Seq[S]) (T, error) {
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
