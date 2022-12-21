package fnio

import (
	"bytes"
	"testing"

	"github.com/kamstrup/fn"
	fntesting "github.com/kamstrup/fn/testing"
)

var text = `hello world
hej verden
hola mundo`

func TestLinesSuite(t *testing.T) {
	createSeq := func() fn.Seq[[]byte] {
		return LinesOf(bytes.NewReader([]byte(text)))
	}
	fntesting.SuiteOf(t, createSeq).Is(
		[]byte("hello world"), []byte("hej verden"), []byte("hola mundo"))
}
