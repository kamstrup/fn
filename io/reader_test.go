package seqio

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

var _ io.Reader = errReader{}
var readError = errors.New("read error")

type errReader struct{}

func (er errReader) Read(p []byte) (n int, err error) {
	return 0, readError
}

func TestReaderSuite(t *testing.T) {
	createSeq := func() BufferSeq {
		return ReaderOf(bytes.NewReader([]byte(text)), make([]byte, 16))
	}

	// text is 33 bytes long
	data := []byte(text)
	data1 := data[:16]
	data2 := data[16:32]
	data3 := data[32:33]

	fntesting.SuiteOf(t, createSeq).Is(
		data1, data2, data3)
}

func TestReaderError(t *testing.T) {
	r := ReaderOf(errReader{}, nil)

	if err := seq.Error(r); err != nil {
		t.Fatal("we should not see an error before we read", err)
	}

	opt, tail := r.First()
	if err := seq.Error(opt); err != readError {
		t.Fatal("opt result must be a read error", err)
	}
	if err := seq.Error(tail); err != readError {
		t.Fatal("tail result must be a read error", err)
	}
}
