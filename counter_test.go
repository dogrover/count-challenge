package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWords(t *testing.T) {
	expectedWords := []Word{"Lorem", "ipsum", "sit", "amet"}
	cases := []struct {
		name string
		data string
		want []Word
	}{
		{"simpleWords", "Lorem ipsum sit amet", expectedWords},
		{"withNewline", "Lorem\nipsum\nsit\namet", expectedWords},
		{"extraSpace", "    Lorem  \n  ipsum\n\n\nsit amet     ", expectedWords},
	}
	// Table-driven tests. Simple loop.
	for _, test := range cases {
		data := strings.NewReader(test.data)
		words := make([]Word, 0, len(test.want))
		for w := range readWords(data) {
			words = append(words, w)
		}
		assert.ElementsMatch(t, test.want, words, "%v: Elements should match", test.name)
	}
}

func TestWordsToTokens(t *testing.T) {
	cases := []struct {
		name string
		data []Word
		want []Token
	}{
		{"toLower", []Word{"Lorem", "IPSUM", "siT", "aMEt"}, []Token{"lorem", "ipsum", "sit", "amet"}},
		{"trimPunct", []Word{"'Lorem'", "(ipsum", "_sit,,,", "amet?"}, []Token{"lorem", "ipsum", "sit", "amet"}},
		{"keepInternalPunct", []Word{"Lo'rem", "_ip_sum_", "si,t", "am-et"}, []Token{"lo'rem", "ip_sum", "si,t", "am-et"}},
		{"unicodeWords", []Word{"Süsse", "Straße", "世界", "'世界'"}, []Token{"süsse", "straße", "世界", "世界"}},
	}
	// Different table strategy: sub-tests.
	for _, test := range cases {
		t.Run(test.name, func(tc *testing.T) {
			wordChan := make(chan Word)
			tokChan := wordsToTokens(wordChan)
			tokens := make([]Token, 0, len(test.data))
			for _, word := range test.data {
				wordChan <- word
				tokens = append(tokens, <-tokChan)
			}
			assert.ElementsMatchf(tc, test.want, tokens, "Elements should match")
		})
	}
}
