package fn

const LenUnknown = -1

type Seq[T any] interface {
	ForEach(f Func1[T])
	ForEachIndex(f Func2[int, T])
	Len() int
	Array() Array[T]
	Take(int) (Array[T], Seq[T])
	TakeWhile(predicate Predicate[T]) (Array[T], Seq[T])
	Skip(int) Seq[T]
	Where(Predicate[T]) Seq[T]
	First() (Opt[T], Seq[T])
}

// TODO:
// Zip
// Assoc(map[K]V)
// Concat(Seq[T], Seq[T])
// seq.Any(pred)/All(pred)
// seq.Split(pred) Seq[Seq[T]]
// Error handling: fn.Must()
// Seq over a channel
// Select on channel
// ctor SeqSlice, SeqChan, SeqSource, SeqRange(int, int), SeqAssoc(map[S]T) Tuple[S,T], AssocBy(seq, FuncMap), SeqOf(tt ... T) Seq[T]
// Reverse?
// testing utils assert/require?
// sorting
// seq.Go(n, f) (n goroutines). Auto-wait, or SeqGo.Wait()?
