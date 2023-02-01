package slice

import (
	"reflect"
	"testing"
)

func TestUniq(t *testing.T) {
	nums := Uniq([]string{"one", "two", "one", "three"})
	if !reflect.DeepEqual(nums, []string{"one", "two", "three"}) {
		t.Fatalf("bad result: %v", nums)
	}
}

func TestDictNoResult(t *testing.T) {
	var (
		corpus       []string
		dict         map[string]uint16
		indexedWords []uint16
	)
	corpus, dict, indexedWords = Dict([]string{"one", "two", "one", "three"}, nil, dict, nil)

	// Build from an empty dict and corpus
	if indexedWords != nil {
		t.Fatalf("indexed words should not be returned")
	}
	if !reflect.DeepEqual(corpus, []string{"one", "two", "three"}) {
		t.Fatalf("bad corpus result: %v", corpus)
	}
	if !reflect.DeepEqual(dict, map[string]uint16{
		"one": 0, "two": 1, "three": 2}) {
		t.Fatalf("bad dict result: %v", dict)
	}

	// Continue building on existing dict and corpus, but add nothing new.
	// Nothing should change
	corpus, dict, _ = Dict([]string{"one", "two"}, corpus, dict, nil)
	if !reflect.DeepEqual(corpus, []string{"one", "two", "three"}) {
		t.Fatalf("bad corpus result: %v", corpus)
	}
	if !reflect.DeepEqual(dict, map[string]uint16{
		"one": 0, "two": 1, "three": 2}) {
		t.Fatalf("bad dict result: %v", dict)
	}

	// Add an extra word to existing dict
	corpus, dict, _ = Dict([]string{"four"}, corpus, dict, nil)
	if !reflect.DeepEqual(corpus, []string{"one", "two", "three", "four"}) {
		t.Fatalf("bad corpus result: %v", corpus)
	}
	if !reflect.DeepEqual(dict, map[string]uint16{
		"one": 0, "two": 1, "three": 2, "four": 3}) {
		t.Fatalf("bad dict result: %v", dict)
	}
}

func TestDictWithResult(t *testing.T) {
	var (
		corpus       []string
		dict         map[string]uint16
		indexedWords = []uint16{}
	)
	corpus, dict, indexedWords = Dict([]string{"one", "two", "one", "three"}, nil, dict, indexedWords)

	if !reflect.DeepEqual(indexedWords, []uint16{0, 1, 0, 2}) {
		t.Fatalf("bad indexed words: %v", indexedWords)
	}
	if !reflect.DeepEqual(corpus, []string{"one", "two", "three"}) {
		t.Fatalf("bad corpus result: %v", corpus)
	}
	if !reflect.DeepEqual(dict, map[string]uint16{
		"one": 0, "two": 1, "three": 2}) {
		t.Fatalf("bad dict result: %v", dict)
	}

	// Continue building on existing dict and corpus, but add nothing new.
	// Nothing should change
	corpus, dict, indexedWords = Dict([]string{"one", "two"}, corpus, dict, indexedWords)
	if !reflect.DeepEqual(indexedWords, []uint16{0, 1}) {
		t.Fatalf("bad indexed words: %v", indexedWords)
	}
	if !reflect.DeepEqual(corpus, []string{"one", "two", "three"}) {
		t.Fatalf("bad corpus result: %v", corpus)
	}
	if !reflect.DeepEqual(dict, map[string]uint16{
		"one": 0, "two": 1, "three": 2}) {
		t.Fatalf("bad dict result: %v", dict)
	}

	// Add an extra word to existing dict
	corpus, dict, indexedWords = Dict([]string{"four"}, corpus, dict, indexedWords)
	if !reflect.DeepEqual(indexedWords, []uint16{3}) {
		t.Fatalf("bad indexed words: %v", indexedWords)
	}
	if !reflect.DeepEqual(corpus, []string{"one", "two", "three", "four"}) {
		t.Fatalf("bad corpus result: %v", corpus)
	}
	if !reflect.DeepEqual(dict, map[string]uint16{
		"one": 0, "two": 1, "three": 2, "four": 3}) {
		t.Fatalf("bad dict result: %v", dict)
	}
}
