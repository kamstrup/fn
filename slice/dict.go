package slice

import (
	"math/bits"

	"github.com/kamstrup/fn/constraints"
)

// Uniq returns a slice with all the unique values from the input slice, in the order they appear.
func Uniq[T comparable](slice []T) []T {
	sz := simpleLog(len(slice))
	dict := make(map[T]struct{}, sz)
	var result []T
	for _, t := range slice {
		if _, seen := dict[t]; !seen {
			dict[t] = struct{}{}
			result = append(result, t)
		}
	}

	return result
}

// Dict is versatile function for deduplicating values, composing indexes, and converting
// values into integer offsets into a corpus.
// It looks at a slice of "words" and builds a corpus of unique words, and a dictionary
// mapping single words to offsets in the corpus.
//
// You can build up state with repeated calls by passing in pre-populated start-corpus and dictionary.
// Most often from results of previous calls to Dict.
//
// If wordIndexBuf is non-nil it will be used to build the indexedWords that maps the "words" argument to the corpus.
// If wordIndexBuf is too small a new buffer will be allocated. If wordIndexBuf is nil then the returned indexedWords
// will also be nil.
//
// # Example 1: Building a Corpus of Unique Words
//
//	words := []string{"one", "two", "one", "three"}
//	corpus, dict, _ := Dict(words, nil, map[string]uint16{}, nil)
//
//	// corpus is now:
//	//    []string{"one", "two", "three"}
//	// dict maps word indexes in the corpus:
//	//    map[string]uint16{"one": 0, "two": 1, "three": 2}
//
//	// We can further build our corpus and dict:
//	moreWords := []string{"four", "one"}
//	corpus, dict, _ = Dict(moreWords, corpus, dict, nil)
//
//	// corpus is now:
//	//    []string{"one", "two", "three", "four"}
//	// dict maps word indexes in the corpus:
//	//    map[string]uint16{"one": 0, "two": 1, "three": 2, "four": 3}
//
// # Example 2: Representing Words as Offsets Reduce A Corpus
//
//		// With the 'corpus' and 'dict' variables from Example 1
//	 lastWords := []string{"one", "two", "five"}
//	 indexedWords := []uint16{} // create a non-nil buffer for the indexed words
//	 corpus, dict, indexedWords = Dict(lastWords, corpus, dict, indexedWords)
//
//	 // indexedWords is now:
//		//    []uint16{0, 1, 4}
//	 // which corresponds exactly with the word-offsets in the updated corpus:
//	 //    []string{"one", "two", "three", "four", "five"}
func Dict[I constraints.Integer, T comparable](words []T, startCorpus []T, startDict map[T]I, wordIndexBuf []I) (corpus []T, dict map[T]I, indexedWords []I) {
	sz := simpleLog(len(words))

	if startCorpus == nil {
		corpus = make([]T, 0, sz)
	} else {
		corpus = startCorpus
	}

	if startDict == nil {
		dict = make(map[T]I, sz)
	} else {
		dict = startDict
	}

	if len(corpus) > len(dict) {
		panic("invalid corpus for dict")
	}

	if wordIndexBuf != nil {
		if cap(wordIndexBuf) >= len(words) {
			indexedWords = wordIndexBuf[:len(words)] // provided buffer has room
		} else {
			indexedWords = make([]I, len(words))
		}
		indexedWords = make([]I, len(words))
		for wordOffset, t := range words {
			corpusIdx, seen := dict[t]
			if !seen {
				corpusIdx = I(len(corpus))
				dict[t] = corpusIdx
				corpus = append(corpus, t)
			}
			indexedWords[wordOffset] = corpusIdx
		}
	} else {
		// No output word index buffer requested
		for _, t := range words {
			if _, seen := dict[t]; !seen {
				dict[t] = I(len(corpus))
				corpus = append(corpus, t)
			}
		}
	}

	return
}

// simpleLog calculates a simplified log2 of an int.
// We use this to estimate sizes of various dictionaries,
// since the number of unique words in a corpus normally
// grows logarithmically in the number of words.
func simpleLog(corpusSize int) int {
	return bits.Len(uint(corpusSize))
}
