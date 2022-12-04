package fnio

import (
	"bytes"
	"io"
	"os"

	"github.com/kamstrup/fn"
)

type readerSeq struct {
	r   io.Reader
	buf []byte
}

func ReaderOf(r io.Reader, buf []byte) fn.Seq[[]byte] {
	if len(buf) == 0 {
		buf = make([]byte, 4096)
	}
	return readerSeq{
		r:   r,
		buf: buf,
	}
}

func (r readerSeq) ForEach(f fn.Func1[[]byte]) fn.Seq[[]byte] {
	for {
		n, err := r.r.Read(r.buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return fn.ErrorOf[[]byte](err)
		}
		f(r.buf[:n])
	}

	return fn.SeqEmpty[[]byte]()
}

// ForEachIndex on a Reader sequence passes the stream offset, not the iteration index to f.
func (r readerSeq) ForEachIndex(f fn.Func2[int, []byte]) fn.Seq[[]byte] {
	offset := 0
	for {
		n, err := r.r.Read(r.buf)
		if err == io.EOF && n == 0 {
			break
		} else if err != nil {
			return fn.ErrorOf[[]byte](err)
		}
		f(offset, r.buf[:n])
		offset += n
	}

	return fn.SeqEmpty[[]byte]()
}

// Len on a Reader is unknown, unless the underlying io.Reader is an *os.File
// in which case it reports the file length, or if it is a something with a Len() int method,
// like a bytes.Buffer.
func (r readerSeq) Len() (int, bool) {
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

// Array redas the full stream and returns it in a byte array.
func (r readerSeq) Array() fn.Array[[]byte] {
	res := r.prepBuffer(0)

	// FIXME: error reporting!?
	fn.Into[[]byte, *bytes.Buffer](res, fn.ByteBuffer, r)
	return fn.ArrayOfArgs(res.Bytes())
}

func (r readerSeq) Take(n int) (fn.Array[[]byte], fn.Seq[[]byte]) {
	res := r.prepBuffer(n)

	for n > 0 {
		reqSz := n
		if reqSz > len(r.buf) {
			reqSz = len(r.buf)
		}

		numRead, err := r.r.Read(r.buf[:reqSz])
		if err == io.EOF && numRead == 0 {
			return fn.ArrayOfArgs(res.Bytes()), fn.SeqEmpty[[]byte]()
		} else if err != nil {
			return fn.ArrayOfArgs(res.Bytes()), fn.ErrorOf[[]byte](err)
		}

		_, _ = res.Write(r.buf[:numRead])
		n -= numRead
	}

	return fn.ArrayOfArgs(res.Bytes()), fn.SeqEmpty[[]byte]()
}

func (r readerSeq) TakeWhile(pred fn.Predicate[[]byte]) (fn.Array[[]byte], fn.Seq[[]byte]) {
	res := &bytes.Buffer{} // we can't pre-alloc -- we really don't know what size buffer we need, could be 0!

	for {
		numRead, err := r.r.Read(r.buf)
		if err == io.EOF && numRead == 0 {
			return fn.ArrayOfArgs(res.Bytes()), fn.SeqEmpty[[]byte]()
		} else if err != nil {
			return fn.ArrayOfArgs(res.Bytes()), fn.ErrorOf[[]byte](err)
		}

		if pred(r.buf) {
			_, _ = res.Write(r.buf[:numRead])
		} else {
			// DANGER: We "unread" r.buf here. We end up in a state where the singlet and r share r.buf.
			// This works out as long as the caller does not user r further, since r will not use the buffer
			// before the singlet is exhausted. We might need to copy the r.buf if this is problematic in practice.
			tail := fn.ConcatOf[[]byte](fn.SingletOf(r.buf), r)
			return fn.ArrayOfArgs(res.Bytes()), tail
		}
	}

	return fn.ArrayOfArgs(res.Bytes()), fn.SeqEmpty[[]byte]()
}

func (r readerSeq) Skip(n int) fn.Seq[[]byte] {
	_, err := io.Copy(io.Discard, io.LimitReader(r.r, int64(n)))
	if err != nil {
		return fn.ErrorOf[[]byte](err)
	}
	return r
}

func (r readerSeq) Where(pred fn.Predicate[[]byte]) fn.Seq[[]byte] {
	return fn.WhereOf[[]byte](r, pred)
}

func (r readerSeq) While(pred fn.Predicate[[]byte]) fn.Seq[[]byte] {
	return fn.WhileOf[[]byte](r, pred)
}

func (r readerSeq) First() (fn.Opt[[]byte], fn.Seq[[]byte]) {
	n, err := r.r.Read(r.buf)
	if err == io.EOF {
		return fn.OptOf[[]byte](r.buf[:n]), fn.SeqEmpty[[]byte]()
	} else if err != nil {
		return fn.OptErr[[]byte](err), fn.ErrorOf[[]byte](err)
	}

	return fn.OptOf(r.buf[:n]), r
}

func (r readerSeq) Shape(shaper fn.FuncMap[[]byte, []byte]) fn.Seq[[]byte] {
	return fn.MapOf[[]byte](r, shaper)
}

// prepBuffer prepares a buffer of a given size. Pass n=0 for sensible default.
func (r readerSeq) prepBuffer(n int) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if sz, ok := r.Len(); ok {
		if n > 0 && sz > n {
			sz = n
		}
		buf.Grow(sz)
	} else if n > 0 {
		buf.Grow(n)
	}
	return buf
}
