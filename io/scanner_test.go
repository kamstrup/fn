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

	fst, tail := r.First()
	if err := fst.Error(); err != readError {
		t.Fatal("opt result must be a read error", err)
	}

	fst, _ = tail.First()
	if err := fst.Error(); err != readError {
		t.Fatal("tail result must be a read error", err)
	}
}
