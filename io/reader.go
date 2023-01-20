package fnio

import (
	"bytes"
	"io"
	"os"

	"github.com/kamstrup/fn"
)

type BufferSeq = fn.Seq[[]byte]
type BufferArray = fn.Array[[]byte]

type Reader struct {
	r   io.Reader
	buf []byte
}

func ReaderOf(r io.Reader, buf []byte) BufferSeq {
	if len(buf) == 0 {
		buf = make([]byte, 4096)
	}
	return Reader{
		r:   r,
		buf: buf,
	}
}

func (r Reader) ForEach(f fn.Func1[[]byte]) BufferSeq {
	for {
		n, err := r.r.Read(r.buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return fn.ErrorOf[[]byte](err)
		}

		f(r.buf[:n]) // might or might not steal r.buf, we have to be defensive!
		r.buf = make([]byte, len(r.buf))
	}

	return fn.SeqEmpty[[]byte]()
}

// ForEachIndex on a Reader sequence passes the stream offset, not the iteration index to f.
func (r Reader) ForEachIndex(f fn.Func2[int, []byte]) BufferSeq {
	i := 0
	for {
		n, err := r.r.Read(r.buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return fn.ErrorOf[[]byte](err)
		}
		f(i, r.buf[:n]) // might or might not steal r.buf, we have to be defensive!
		r.buf = make([]byte, len(r.buf))
		i++
	}

	return fn.SeqEmpty[[]byte]()
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
	return fn.LenUnknown, false
}

// ByteLen on a Reader is unknown, unless the underlying io.Reader is an *os.File
// in which case it reports the file length, or if it is a something with a Len() int method,
// like a bytes.Buffer.
func (r Reader) ByteLen() (int, bool) {
	// Check if reader is a File
	if f, ok := r.r.(*os.File); ok {
		if st, err := f.Stat(); err != nil {
			return fn.LenUnknown, false
		} else {
			return int(st.Size()), true
		}
	}

	// Check if reader is a bytes.Buffer or something else with a Len()
	if withLen, ok := r.r.(interface{ Len() int }); ok {
		return withLen.Len(), true
	}

	return fn.LenUnknown, false
}

func (r Reader) Array() BufferArray {
	// TODO: if size is well-defined: alloc 1 continuous stride and do 1 read call, and sub-divide into buffers via slicing

	return fn.Into(nil, fn.Append[[]byte], r.seq()).Or(nil) // careful: errors silently dropped
}

func (r Reader) Take(n int) (BufferArray, BufferSeq) {
	var (
		res  [][]byte
		tail = r
	)

	for i := 0; i < n; i++ {
		numRead, err := r.r.Read(r.buf)
		if err == io.EOF && numRead == 0 {
			return res, fn.SeqEmpty[[]byte]()
		} else if err != nil {
			return res, fn.ErrorOf[[]byte](err)
		}

		res = append(res, r.buf[:numRead])
		r.buf = make([]byte, len(r.buf))
	}

	return res, tail
}

func (r Reader) TakeWhile(pred fn.Predicate[[]byte]) (BufferArray, BufferSeq) {
	var res [][]byte // we can't pre-alloc -- we really don't know what size buffer we need, could be 0!

	for {
		numRead, err := r.r.Read(r.buf)
		if err == io.EOF && numRead == 0 {
			return res, fn.SeqEmpty[[]byte]()
		} else if err != nil {
			return res, fn.ErrorOf[[]byte](err)
		}

		if pred(r.buf[:numRead]) {
			res = append(res, r.buf[:numRead])
			r.buf = make([]byte, len(r.buf))
		} else {
			// DANGER: We "unread" r.buf here. We end up in a state where the singlet and r share r.buf.
			// This works out as long as the caller does not user r further, since r will not use the buffer
			// before the singlet is exhausted. We might need to copy the r.buf if this is problematic in practice.
			tail := fn.ConcatOf[[]byte](fn.SingletOf(r.buf), r)
			return res, tail
		}
	}

	return res, fn.SeqEmpty[[]byte]()
}

func (r Reader) Skip(n int) BufferSeq {
	// Skip n buffers
	_, err := io.Copy(io.Discard, io.LimitReader(r.r, int64(n*len(r.buf))))
	if err != nil {
		return fn.ErrorOf[[]byte](err)
	}
	return r
}

func (r Reader) Where(pred fn.Predicate[[]byte]) BufferSeq {
	return fn.WhereOf[[]byte](r, pred)
}

func (r Reader) While(pred fn.Predicate[[]byte]) BufferSeq {
	return fn.WhileOf[[]byte](r, pred)
}

func (r Reader) First() (fn.Opt[[]byte], BufferSeq) {
	n, err := r.r.Read(r.buf)
	if err == io.EOF {
		if n == 0 {
			return fn.OptEmpty[[]byte](), fn.SeqEmpty[[]byte]()
		}
		return fn.OptOf[[]byte](r.buf[:n]), fn.SeqEmpty[[]byte]()
	} else if err != nil {
		return fn.OptErr[[]byte](err), fn.ErrorOf[[]byte](err)
	}

	return fn.OptOf(r.buf[:n]), r
}

func (r Reader) Map(shaper fn.FuncMap[[]byte, []byte]) BufferSeq {
	return fn.MapOf[[]byte](r, shaper)
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
