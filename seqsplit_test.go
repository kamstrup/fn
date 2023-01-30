package fn_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kamstrup/fn"
	fntesting "github.com/kamstrup/fn/testing"
)

func compareSplits(t *testing.T, seq fn.Seq[fn.Seq[byte]], words ...string) {
	t.Helper()

	arr := seq.Values()
	if len(arr) != len(words) {
		t.Errorf("expected length %d, got %d", len(words), len(words))
	}

	for i, w := range words {
		resultWord := string(arr[i].Values())
		if resultWord != w {
			t.Errorf("Expected %q, got %q, at index %d", w, resultWord, i)
		}
	}
}

func TestSplit(t *testing.T) {
	words := fn.SplitOf(fn.StringOf("hello.world"), func(c byte) fn.SplitChoice {
		if c == '.' {
			return fn.SplitSeparate
		}
		return fn.SplitKeep
	})

	compareSplits(t, words, "hello", "world")
}

func TestSplitFlatten(t *testing.T) {
	// Split by each letter
	letters := fn.SplitOf(fn.StringOf("hello"), func(c byte) fn.SplitChoice {
		return fn.SplitSeparateKeep
	})

	words := fn.FlattenOf(letters)
	fntesting.TestOf(t, words).Is([]byte("hello")...)
}

func TestSplitSuite(t *testing.T) {
	createSeq := func() fn.Seq[fn.Seq[byte]] {
		// Split by each letter
		return fn.SplitOf(fn.StringOf("hello"), func(c byte) fn.SplitChoice {
			return fn.SplitSeparateKeep
		})
	}

	fntesting.SuiteOf(t, createSeq).
		WithEqual(func(s1, s2 fn.Seq[byte]) bool {
			return reflect.DeepEqual(s1.Values(), s2.Values())
		}).
		Is(fn.SingletOf[byte]('h'), fn.SingletOf[byte]('e'), fn.SingletOf[byte]('l'),
			fn.SingletOf[byte]('l'), fn.SingletOf[byte]('o'))
}

func TestSplitError(t *testing.T) {
	theError := errors.New("the error")
	split := fn.SplitOf(fn.ErrorOf[int](theError), func(_ int) fn.SplitChoice {
		return fn.SplitKeep
	})

	if err := fn.Error(split); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
