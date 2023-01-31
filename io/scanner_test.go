package seqio

import (
	"bytes"
	"testing"

	"github.com/kamstrup/fn/seq"
	fntesting "github.com/kamstrup/fn/testing"
)

var text = `hello world
hej verden
hola mundo`

func TestLinesSuite(t *testing.T) {
	createSeq := func() seq.Seq[[]byte] {
		return LinesOf(bytes.NewReader([]byte(text)))
	}
	fntesting.SuiteOf(t, createSeq).Is(
		[]byte("hello world"), []byte("hej verden"), []byte("hola mundo"))
}

func TestLinesError(t *testing.T) {
	r := LinesOf(errReader{})

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
