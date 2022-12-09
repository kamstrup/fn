package fn_test

import (
	"testing"

	"github.com/kamstrup/fn"
	"github.com/kamstrup/fn/testing"
)

func TestString(t *testing.T) {
	createSeq := func() fn.Seq[byte] {
		return fn.StringOf("hello world")
	}
	fntesting.SuiteOf[byte](t, createSeq).Is([]byte("hello world")...)
}
