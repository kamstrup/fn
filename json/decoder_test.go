package seqjson

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/kamstrup/fn/seq"
	fntesting "github.com/kamstrup/fn/testing"
)

var twoPayloads = `{"Foo": 27}
{"Foo": 28}`

var errPayload = `{"Foo": 27}
{"Foo": `

type Payload struct {
	Foo int
}

func TestDecoderOk(t *testing.T) {
	createSeq := func() seq.Seq[Payload] {
		rdr := bytes.NewReader([]byte(twoPayloads))
		dec := json.NewDecoder(rdr)
		return DecoderOf[Payload](dec)
	}

	fntesting.SuiteOf(t, createSeq).Is(Payload{Foo: 27}, Payload{Foo: 28})
}

func TestDecoderError(t *testing.T) {
	createSeq := func() seq.Seq[Payload] {
		rdr := bytes.NewReader([]byte(errPayload))
		dec := json.NewDecoder(rdr)
		return DecoderOf[Payload](dec)
	}

	// We should read the first element
	elems := createSeq().ToSlice()
	if !reflect.DeepEqual(elems.AsSlice(), []Payload{{Foo: 27}}) {
		t.Fatalf("bad result: %v", elems)
	}

	// We should read the first element, and then ForEach returns an error
	elems = []Payload{}
	resultSeq := createSeq().ForEach(func(p Payload) {
		elems = append(elems, p)
	})
	if !reflect.DeepEqual(elems.AsSlice(), []Payload{{Foo: 27}}) {
		t.Fatalf("bad result: %v", elems)
	}
	if err := seq.Error(resultSeq); err != io.ErrUnexpectedEOF {
		t.Fatalf("expected EOF error, got: %v", err)
	}
}
