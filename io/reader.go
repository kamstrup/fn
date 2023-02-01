package seqio

import (
	"bytes"
	"io"
	"os"

	"github.com/kamstrup/fn/opt"
	"github.com/kamstrup/fn/seq"
)

type BufferSeq = seq.Seq[[]byte]
type BufferArray = seq.Slice[[]byte]

type Reader struct {
	r   io.Reader
	buf []byte
}

// ReaderOf creates a stateful seq.Seq wrapping an io.Reader.
// A buffer to use may optionally be provided.
//
// Errors can be detected be using seq.Error() on seqs and opts
// returned from the reader's methods.
func ReaderOf(r io.Reader, buf []byte) BufferSeq {
	if len(buf) == 0 {
		buf = make([]byte, 4096)
	}
	return Reader{
		r:   r,
		buf: buf,
	}
}

func (r Reader) ForEach(f seq.Func1[[]byte]) BufferSeq {
	for {
		n, err := r.r.Read(r.buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return seq.ErrorOf[[]byte](err)
		}

		f(r.buf[:n]) // might or might not steal r.buf, we have to be defensive!
		r.buf = make([]byte, len(r.buf))
	}

	return seq.Empty[[]byte]()
}

// ForEachIndex on a Reader sequence passes the stream offset, not the iteration index to f.
func (r Reader) ForEachIndex(f seq.Func2[int, []byte]) BufferSeq {
	i := 0
	for {
		n, err := r.r.Read(r.buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return seq.ErrorOf[[]byte](err)
		}
		f(i, r.buf[:n]) // might or might not steal r.buf, we have to be defensive!
		r.buf = make([]byte, len(r.buf))
		i++
	}

	return seq.Empty[[]byte]()
}

// Len on a Reader is unknown, unless the underlying io.Reader is an *os.File
// or something with a Len() int method, in which case it returns the number of buffers
// this Seq will produce when executed.
func (r Reader) Len() (int, bool) {
	if sz, ok := r.ByteLen(); ok {
		if rem := sz % len(r.buf); rem != 0 {
			return sz/len(r.buf) + 1, true
		}
		return sz / len(r.buf), true
	}
	return seq.LenUnknown, false
}

// ByteLen on a Reader is unknown, unless the underlying io.Reader is an *os.File
// or if it is a something with a Len() int method, like a bytes.Buffer.
func (r Reader) ByteLen() (int, bool) {
	// Check if reader is a File
	if f, ok := r.r.(*os.File); ok {
		if st, err := f.Stat(); err != nil {
			return seq.LenUnknown, false
		} else {
			return int(st.Size()), true
		}
	}

	// Check if reader is a bytes.Buffer or something else with a Len()
	if withLen, ok := r.r.(interface{ Len() int }); ok {
		return withLen.Len(), true
	}

	return seq.LenUnknown, false
}

func (r Reader) ToSlice() BufferArray {
	// TODO: if size is well-defined: alloc 1 continuous stride and do 1 read call, and sub-divide into buffers via slicing

	return seq.Reduce(seq.Append[[]byte], nil, r.seq()).Or(nil) // careful: errors silently dropped
}

func (r Reader) Take(n int) (BufferArray, BufferSeq) {
	var (
		res  [][]byte
		tail = r
	)

	for i := 0; i < n; i++ {
		numRead, err := r.r.Read(r.buf)
		if err == io.EOF && numRead == 0 {
			return res, seq.Empty[[]byte]()
		} else if err != nil {
			return res, seq.ErrorOf[[]byte](err)
		}

		res = append(res, r.buf[:numRead])
		r.buf = make([]byte, len(r.buf))
	}

	return res, tail
}

func (r Reader) TakeWhile(pred seq.Predicate[[]byte]) (BufferArray, BufferSeq) {
	var res [][]byte // we can't pre-alloc -- we really don't know what size buffer we need, could be 0!

	for {
		numRead, err := r.r.Read(r.buf)
		if err == io.EOF && numRead == 0 {
			return res, seq.Empty[[]byte]()
		} else if err != nil {
			return res, seq.ErrorOf[[]byte](err)
		}

		if pred(r.buf[:numRead]) {
			res = append(res, r.buf[:numRead])
			r.buf = make([]byte, len(r.buf))
		} else {
			// DANGER: We "unread" r.buf here. We end up in a state where the singlet and r share r.buf.
			// This works out as long as the caller does not user r further, since r will not use the buffer
			// before the singlet is exhausted. We might need to copy the r.buf if this is problematic in practice.
			tail := seq.PrependOf[[]byte](r.buf, r)
			return res, tail
		}
	}

	return res, seq.Empty[[]byte]()
}

func (r Reader) Skip(n int) BufferSeq {
	// Skip n buffers
	_, err := io.Copy(io.Discard, io.LimitReader(r.r, int64(n*len(r.buf))))
	if err != nil {
		return seq.ErrorOf[[]byte](err)
	}
	return r
}

func (r Reader) Where(pred seq.Predicate[[]byte]) BufferSeq {
	return seq.WhereOf[[]byte](r, pred)
}

func (r Reader) While(pred seq.Predicate[[]byte]) BufferSeq {
	return seq.WhileOf[[]byte](r, pred)
}

func (r Reader) First() (opt.Opt[[]byte], BufferSeq) {
	n, err := r.r.Read(r.buf)
	if err == io.EOF {
		if n == 0 {
			return opt.Empty[[]byte](), seq.Empty[[]byte]()
		}
		return opt.Of[[]byte](r.buf[:n]), seq.Empty[[]byte]()
	} else if err != nil {
		return opt.ErrorOf[[]byte](err), seq.ErrorOf[[]byte](err)
	}

	return opt.Of(r.buf[:n]), r
}

func (r Reader) Map(shaper seq.FuncMap[[]byte, []byte]) BufferSeq {
	return seq.MappingOf[[]byte](r, shaper)
}

// prepBuffer prepares a buffer of a given size. Pass n=0 for sensible default.
func (r Reader) prepBuffer(n int) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if sz, ok := r.ByteLen(); ok {
		if n > 0 && sz > n {
			sz = n
		}
		buf.Grow(sz)
	} else if n > 0 {
		buf.Grow(n)
	}
	return buf
}

func (r Reader) seq() BufferSeq {
	return r
}
