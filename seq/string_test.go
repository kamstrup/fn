package seq_test

import (
	"testing"

	"github.com/kamstrup/fn/seq"
	"github.com/kamstrup/fn/testing"
)

func TestString(t *testing.T) {
	createSeq := func() seq.Seq[byte] {
		return seq.StringOf("hello world")
	}
	fntesting.SuiteOf[byte](t, createSeq).Is([]byte("hello world")...)
}
