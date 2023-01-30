package fn

import "github.com/kamstrup/fn/opt"

type String string

func StringOf(s string) Seq[byte] {
	return String(s)
}

func StringAs(s string) String {
	return String(s)
}

func (s String) ForEach(f Func1[byte]) Seq[byte] {
	// a for-range loop on a string iterates codepoints, and not bytes,
	// so we do a handrolled for-loop to go byte by byte
	for i := 0; i < len(s); i++ {
		f(s[i])
	}

	return SeqEmpty[byte]()
}

func (s String) ForEachIndex(f Func2[int, byte]) Seq[byte] {
	// a for-range loop on a string iterates codepoints, and not bytes,
	// so we do a handrolled for-loop to go byte by byte
	for i := 0; i < len(s); i++ {
		f(i, s[i])
	}

	return SeqEmpty[byte]()
}

func (s String) Len() (int, bool) {
	return len(s), true
}

func (s String) Array() Array[byte] {
	return []byte(s)
}

func (s String) Take(n int) (Array[byte], Seq[byte]) {
	if len(s) <= n {
		return []byte(s), SeqEmpty[byte]()
	}
	return []byte(s[:n]), s[n:]
}

func (s String) TakeWhile(pred Predicate[byte]) (Array[byte], Seq[byte]) {
	for i := 0; i < len(s); i++ {
		if !pred(s[i]) {
			return []byte(s[:i]), s[i:]
		}
	}
	return []byte(s), SeqEmpty[byte]()
}

func (s String) Skip(n int) Seq[byte] {
	if len(s) <= n {
		return SeqEmpty[byte]()
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
