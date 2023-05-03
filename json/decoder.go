package seqjson

import (
	"encoding/json"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

type decoderSeq[T any] struct {
	dec *json.Decoder
}

func DecoderOf[T any](dec *json.Decoder) seq.Seq[T] {
	return decoderSeq[T]{
		dec: dec,
	}
}

func (d decoderSeq[T]) ForEach(f seq.Func1[T]) opt.Opt[T] {
	for d.dec.More() {
		var t T
		if err := d.dec.Decode(&t); err != nil {
			return opt.ErrorOf[T](err)
		}
		f(t)
	}
	return opt.Zero[T]()
}

func (d decoderSeq[T]) ForEachIndex(f seq.Func2[int, T]) opt.Opt[T] {
	i := 0
	for d.dec.More() {
		var t T
		if err := d.dec.Decode(&t); err != nil {
			return opt.ErrorOf[T](err)
		}
		f(i, t)
		i++
	}
	return opt.Zero[T]()
}

func (d decoderSeq[T]) Len() (int, bool) {
	return seq.LenUnknown, false
}

func (d decoderSeq[T]) ToSlice() seq.Slice[T] {
	var arr []T
	d.ForEach(func(t T) {
		arr = append(arr, t)
	})
	return arr
}

func (d decoderSeq[T]) Limit(n int) seq.Seq[T] {
	return seq.LimitOf[T](d, n)
}

func (d decoderSeq[T]) Take(n int) (seq.Slice[T], seq.Seq[T]) {
	var (
		arr = make([]T, 0, n)
		i   = 0
	)
	for d.dec.More() && i < n {
		var t T
		if err := d.dec.Decode(&t); err != nil {
			return arr, seq.ErrorOf[T](err)
		}
		arr = append(arr, t)
		i++
	}
	return arr, seq.Empty[T]()
}

func (d decoderSeq[T]) TakeWhile(pred seq.Predicate[T]) (seq.Slice[T], seq.Seq[T]) {
	var arr []T
	for d.dec.More() {
		var t T
		if err := d.dec.Decode(&t); err != nil {
			return arr, seq.ErrorOf[T](err)
		}

		if pred(t) {
			arr = append(arr, t)
		} else {
			return arr, seq.PrependOf[T](t, d)
		}
	}
	return arr, seq.Empty[T]()
}

func (d decoderSeq[T]) Skip(n int) seq.Seq[T] {
	i := 0
	for d.dec.More() && i < n {
		var t T
		if err := d.dec.Decode(&t); err != nil {
			return seq.ErrorOf[T](err)
		}
		i++
	}
	return d
}

func (d decoderSeq[T]) Where(pred seq.Predicate[T]) seq.Seq[T] {
	return seq.WhereOf[T](d, pred)
}

func (d decoderSeq[T]) While(pred seq.Predicate[T]) seq.Seq[T] {
	return seq.WhileOf[T](d, pred)
}

func (d decoderSeq[T]) First() (opt.Opt[T], seq.Seq[T]) {
	if d.dec.More() {
		var t T
		if err := d.dec.Decode(&t); err != nil {
			return opt.ErrorOf[T](err), seq.ErrorOf[T](err)
		}
		return opt.Of(t), d
	}

	return opt.Empty[T](), seq.Empty[T]()
}

func (d decoderSeq[T]) Map(f seq.FuncMap[T, T]) seq.Seq[T] {
	return seq.MappingOf[T, T](d, f)
}
