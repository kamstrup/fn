package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
)

func TestString(t *testing.T) {
	createSeq := func() fn.Seq[byte] {
		return fn.StringOf("hello world")
	}
	fn.SeqTestSuite[byte](t, createSeq).Is([]byte("hello world")...)
}
