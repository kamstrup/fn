package fnio

import (
	"bufio"
	"io"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

type scannerSeq struct {
	scanner *bufio.Scanner
}

// ScannerOf creates a stateful seq splits a io.Reader into delimited chunks.
// If you want to split a reader into lines you can use LinesOf.
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

// LinesOf returns a stateful seq splitting an io.Reader into lines using a standard bufio.Scanner.
func LinesOf(r io.Reader) BufferSeq {
	return scannerSeq{
		scanner: bufio.NewScanner(r),
	}
}

func (s scannerSeq) ForEach(f seq.Func1[[]byte]) BufferSeq {
	for s.scanner.Scan() {
		f(s.scanner.Bytes())
	}
	return s.errOrEmpty()
}

func (s scannerSeq) ForEachIndex(f seq.Func2[int, []byte]) BufferSeq {
	for i := 0; s.scanner.Scan(); i++ {
		f(i, s.scanner.Bytes())
	}
	return s.errOrEmpty()
}

func (s scannerSeq) Len() (int, bool) {
	return seq.LenUnknown, false
}

func (s scannerSeq) Values() BufferArray {
	tokens := seq.Into(nil, func(tokens [][]byte, tok []byte) [][]byte {
		dupTok := append([]byte{}, tok...) // scanner owns tok, so we copy it
		return append(tokens, dupTok)
	}, s.seq()).Or(nil) // careful: errors silently dropped
	return tokens
}

func (s scannerSeq) Take(n int) (BufferArray, BufferSeq) {
	var tokens [][]byte
	for i := 0; i < n && s.scanner.Scan(); i++ {
		tok := append([]byte{}, s.scanner.Bytes()...) // scanner owns the Bytes() buffer
		tokens = append(tokens, tok)
	}

	return tokens, s
}

func (s scannerSeq) TakeWhile(pred seq.Predicate[[]byte]) (BufferArray, BufferSeq) {
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
		return tokens, seq.ConcatOf(seq.SingletOf(tok), BufferSeq(s))
	}

	return tokens, s
}

func (s scannerSeq) Skip(n int) BufferSeq {
	for i := 0; i < n && s.scanner.Scan(); i++ {
		// skip
	}

	return s
}

func (s scannerSeq) Where(pred seq.Predicate[[]byte]) BufferSeq {
	return seq.WhereOf[[]byte](s, pred)
}

func (s scannerSeq) While(pred seq.Predicate[[]byte]) BufferSeq {
	return seq.WhileOf[[]byte](s, pred)
}

func (s scannerSeq) First() (opt.Opt[[]byte], BufferSeq) {
	if s.scanner.Scan() {
		tok := append([]byte{}, s.scanner.Bytes()...) // scanner owns Bytes() buffer
		return opt.Of(tok), s
	}
	if err := s.scanner.Err(); err != nil {
		return opt.ErrorOf[[]byte](err), seq.ErrorOf[[]byte](err)
	}
	return opt.Empty[[]byte](), s
}

func (s scannerSeq) Map(m seq.FuncMap[[]byte, []byte]) BufferSeq {
	return seq.MapOf[[]byte, []byte](s, m)
}

// Error implements the contract for the seq.Error function.
func (s scannerSeq) Error() error {
	return s.scanner.Err()
}

func (s scannerSeq) errOrEmpty() BufferSeq {
	if err := s.Error(); err != nil {
		return seq.ErrorOf[[]byte](err)
	}
	return seq.SeqEmpty[[]byte]()
}

func (s scannerSeq) seq() BufferSeq {
	return s
}
