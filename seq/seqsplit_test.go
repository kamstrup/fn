package seq_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kamstrup/fn/seq"
	fntesting "github.com/kamstrup/fn/testing"
)

func compareSplits(t *testing.T, seq seq.Seq[seq.Seq[byte]], words ...string) {
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
	words := seq.SplitOf(seq.StringOf("hello.world"), func(c byte) seq.SplitChoice {
		if c == '.' {
			return seq.SplitSeparate
		}
		return seq.SplitKeep
	})

	compareSplits(t, words, "hello", "world")
}

func TestSplitFlatten(t *testing.T) {
	// Split by each letter
	letters := seq.SplitOf(seq.StringOf("hello"), func(c byte) seq.SplitChoice {
		return seq.SplitSeparateKeep
	})

	words := seq.FlattenOf(letters)
	fntesting.TestOf(t, words).Is([]byte("hello")...)
}

func TestSplitSuite(t *testing.T) {
	createSeq := func() seq.Seq[seq.Seq[byte]] {
		// Split by each letter
		return seq.SplitOf(seq.StringOf("hello"), func(c byte) seq.SplitChoice {
			return seq.SplitSeparateKeep
		})
	}

	fntesting.SuiteOf(t, createSeq).
		WithEqual(func(s1, s2 seq.Seq[byte]) bool {
			return reflect.DeepEqual(s1.Values(), s2.Values())
		}).
		Is(seq.SingletOf[byte]('h'), seq.SingletOf[byte]('e'), seq.SingletOf[byte]('l'),
			seq.SingletOf[byte]('l'), seq.SingletOf[byte]('o'))
}

func TestSplitError(t *testing.T) {
	theError := errors.New("the error")
	split := seq.SplitOf(seq.ErrorOf[int](theError), func(_ int) seq.SplitChoice {
		return seq.SplitKeep
	})

	if err := seq.Error(split); err != theError {
		t.Fatalf("Expected 'the error', found: %s", err)
	}
}
