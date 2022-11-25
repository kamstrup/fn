package fn

type String string

func StringOf(s string) Seq[byte] {
	return String(s)
}

func (s String) ForEach(f Func1[byte]) {
	// a for-range loop on a string iterates codepoints, and not bytes,
	// so we do a handrolled for-loop to go byte by byte
	for i := 0; i < len(s); i++ {
		f(s[i])
	}
}

func (s String) ForEachIndex(f Func2[int, byte]) {
	// a for-range loop on a string iterates codepoints, and not bytes,
	// so we do a handrolled for-loop to go byte by byte
	for i := 0; i < len(s); i++ {
		f(i, s[i])
	}
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

func (s String) First() (Opt[byte], Seq[byte]) {
	if len(s) == 0 {
		return OptEmpty[byte](), s
	}
	return OptOf(s[0]), s[1:]
}
