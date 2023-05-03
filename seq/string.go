package seq

import (
	"strings"

	"github.com/kamstrup/fn/opt"
)

// String is a type wrapper for standard go strings.
// Any builtin or go syntax that applies to a string also applies to String.
// You may fx. call len() and do for-range loops on seq.String.
type String string

// StringOf returns the string interpreted as a seq of bytes.
func StringOf(s string) Seq[byte] {
	return String(s)
}

func StringAs(s string) String {
	return String(s)
}

func (s String) ForEach(f Func1[byte]) opt.Opt[byte] {
	// a for-range loop on a string iterates codepoints, and not bytes,
	// so we do a handrolled for-loop to go byte by byte
	for i := 0; i < len(s); i++ {
		f(s[i])
	}

	return opt.Zero[byte]()
}

func (s String) ForEachIndex(f Func2[int, byte]) opt.Opt[byte] {
	// a for-range loop on a string iterates codepoints, and not bytes,
	// so we do a handrolled for-loop to go byte by byte
	for i := 0; i < len(s); i++ {
		f(i, s[i])
	}

	return opt.Zero[byte]()
}

func (s String) Len() (int, bool) {
	return len(s), true
}

func (s String) ToSlice() Slice[byte] {
	return []byte(s)
}

func (s String) Limit(n int) Seq[byte] {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func (s String) Take(n int) (Slice[byte], Seq[byte]) {
	if len(s) <= n {
		return []byte(s), Empty[byte]()
	}
	return []byte(s[:n]), s[n:]
}

func (s String) TakeWhile(pred Predicate[byte]) (Slice[byte], Seq[byte]) {
	for i := 0; i < len(s); i++ {
		if !pred(s[i]) {
			return []byte(s[:i]), s[i:]
		}
	}
	return []byte(s), Empty[byte]()
}

func (s String) Skip(n int) Seq[byte] {
	if len(s) <= n {
		return Empty[byte]()
	}
	return s[n:]
}

func (s String) Where(pred Predicate[byte]) Seq[byte] {
	return whereSeq[byte]{
		seq:  s,
		pred: pred,
	}
}

func (s String) While(pred Predicate[byte]) Seq[byte] {
	return whileSeq[byte]{
		seq:  s,
		pred: pred,
	}
}

func (s String) First() (opt.Opt[byte], Seq[byte]) {
	if len(s) == 0 {
		return opt.Empty[byte](), s
	}
	return opt.Of(s[0]), s[1:]
}

func (s String) Map(shaper FuncMap[byte, byte]) Seq[byte] {
	return mappedSeq[byte, byte]{
		f:   shaper,
		seq: s,
	}
}

func (s String) HasSuffix(sfx string) bool {
	return strings.HasSuffix(string(s), sfx)
}

func (s String) HasPrefix(pfx string) bool {
	return strings.HasPrefix(string(s), pfx)
}

func (s String) Contains(sub string) bool {
	return strings.Contains(string(s), sub)
}
