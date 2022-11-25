package fn

import "testing"

func TestString(t *testing.T) {
	createSeq := func() Seq[byte] {
		return StringOf("hello world")
	}
	SeqTestSuite[byte](t, createSeq).Is([]byte("hello world")...)
}
