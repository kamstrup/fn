package fnio

import (
	"bufio"
	"io"

	"github.com/kamstrup/fn"
)

type scannerSeq struct {
	scanner *bufio.Scanner
}

func ScannerOf(r io.Reader, split bufio.SplitFunc, buf []byte, maxBuf int) BufferSeq {
	scanner := bufio.NewScanner(r)
	if split != nil {
		scanner.Split(split)
	}
	if buf != nil {
		scanner.Buffer(buf, maxBuf)
	}
	return scannerSeq{scanner: scanner}
}

func LinesOf(r io.Reader) BufferSeq {
	return scannerSeq{
		scanner: bufio.NewScanner(r),
	}
}

func (s scannerSeq) ForEach(f fn.Func1[[]byte]) BufferSeq {
	for s.scanner.Scan() {
		f(s.scanner.Bytes())
	}
	return s.errOrEmpty()
}

func (s scannerSeq) ForEachIndex(f fn.Func2[int, []byte]) BufferSeq {
	for i := 0; s.scanner.Scan(); i++ {
		f(i, s.scanner.Bytes())
	}
	return s.errOrEmpty()
}

func (s scannerSeq) Len() (int, bool) {
	return fn.LenUnknown, false
}

func (s scannerSeq) Array() BufferArray {
	tokens := fn.Into[[]byte, [][]byte](nil, func(tokens [][]byte, tok []byte) [][]byte {
		dupTok := append([]byte{}, tok...) // scanner owns tok, so we copy it
		return append(tokens, dupTok)
	}, s)
	return fn.ArrayAs(tokens)
}

func (s scannerSeq) Take(n int) (BufferArray, BufferSeq) {
	var tokens [][]byte
	for i := 0; i < n && s.scanner.Scan(); i++ {
		tok := append([]byte{}, s.scanner.Bytes()...) // scanner owns the Bytes() buffer
		tokens = append(tokens, tok)
	}

	return tokens, s
}

func (s scannerSeq) TakeWhile(pred fn.Predicate[[]byte]) (BufferArray, BufferSeq) {
	var (
		tokens [][]byte
		tok    []byte
	)
	for s.scanner.Scan() {
		tok = append([]byte{}, s.scanner.Bytes()...) // scanner owns Bytes() buffer
		if pred(tok) {
			tokens = append(tokens, tok)
			tok = nil
		} else {
			break
		}
	}

	if s.Error() != nil {
		return tokens, s // FIXME: we drop the last token on error -- would need a ConcatOfWithError() :-/
	}

	// tok did not match pred, so push it back onto the seq
	if tok != nil {
		return tokens, fn.ConcatOf(fn.SingletOf(tok), BufferSeq(s))
	}

	return tokens, s
}

func (s scannerSeq) Skip(n int) BufferSeq {
	for i := 0; i < n && s.scanner.Scan(); i++ {
		// skip
	}

	return s
}

func (s scannerSeq) Where(pred fn.Predicate[[]byte]) BufferSeq {
	return fn.WhereOf[[]byte](s, pred)
}

func (s scannerSeq) While(pred fn.Predicate[[]byte]) BufferSeq {
	return fn.WhileOf[[]byte](s, pred)
}

func (s scannerSeq) First() (fn.Opt[[]byte], BufferSeq) {
	if s.scanner.Scan() {
		tok := append([]byte{}, s.scanner.Bytes()...) // scanner owns Bytes() buffer
		return fn.OptOf(tok), s
	}
	if err := s.scanner.Err(); err != nil {
		return fn.OptErr[[]byte](err), fn.ErrorOf[[]byte](err)
	}
	return fn.OptEmpty[[]byte](), s
}

func (s scannerSeq) Map(m fn.FuncMap[[]byte, []byte]) BufferSeq {
	return fn.MapOf[[]byte, []byte](s, m)
}

// Error implements the contract for the fn.Error function.
func (s scannerSeq) Error() error {
	return s.scanner.Err()
}

func (s scannerSeq) errOrEmpty() BufferSeq {
	if err := s.Error(); err != nil {
		return fn.ErrorOf[[]byte](err)
	}
	return fn.SeqEmpty[[]byte]()
}
