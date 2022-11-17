package fn

const LenUnknown = -1

type Seq[T comparable] interface {
	ForEach(f Func1[T])
	ForEachIndex(f Func2[int, T])
	Len() int
	Array() Array[T]
	Take(int) (Seq[T], Seq[T])
	First() (Opt[T], Seq[T])
}

// TODO:
// seq.First() (t T, rest Seq[T])
// seq.Take(int) (Array[T], rest Seq[T])
// seq.TakeWhile(pred) (head Array[T], rest Seq[T])
// Where() or Collect() ?? (not Select! use for chan)
// Zip
// Assoc
// Concat
// Opt[T] // Value, Or, Return() -> T??
// Tuple (as Seq)
// seq.Any(pred)/All(pred)
// seq.Split(pred) Seq[Seq[T]]
// Error handling: fn.Must()
// lib name: "fn"
// Seq over a channel
// Select on Channel
// ctor SeqSlice, SeqChan, SeqSource, SeqRange(int, int), SeqAssoc(map[S]T) Tuple[S,T], AssocBy(seq, FuncMap), SeqOf(tt ... T) Seq[T]
// Reverse?
// testing utils assert/require?
