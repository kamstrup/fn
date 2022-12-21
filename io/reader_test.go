package fnio

import (
	"bytes"
	"testing"

	fntesting "github.com/kamstrup/fn/testing"
)

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
