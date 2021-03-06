package main

import (
	"io"
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
	// Table-driven tests with a simple loop.
	for _, test := range cases {
		srcChan := make(chan io.Reader)
		go func() {
			srcChan <- strings.NewReader(test.data)
			close(srcChan)
		}()
		words := make([]Word, 0, len(test.want))
		for w := range readWords(srcChan) {
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
	// Different table strategy: sub-tests. Better failure handling, and test
	// identification! Have to be careful, though, not to mix up the outer
	// testing.T with the inner one ("t" vs. "tc", here)
	for _, test := range cases {
		t.Run(test.name, func(tc *testing.T) {
			wordChan := make(chan Word)
			tokChan := wordsToTokens(wordChan)
			tokens := make([]Token, 0, len(test.data))
			for _, word := range test.data {
				wordChan <- word
				tokens = append(tokens, <-tokChan)
			}
			close(wordChan)
			assert.ElementsMatchf(tc, test.want, tokens, "Elements should match")
		})
	}
}

func TestGetChunks(t *testing.T) {
	noChunks := []Chunk{}
	oneChunk := []Chunk{
		{"lorem", "ipsum", "dolor"},
	}
	multiChunks := []Chunk{
		{"lorem", "ipsum", "dolor"},
		{"ipsum", "dolor", "sit"},
		{"dolor", "sit", "amet"},
	}
	cases := []struct {
		name string
		data []Token
		want []Chunk
	}{
		{"noTokens", []Token{}, noChunks},
		{"oneToken", []Token{"lorem"}, noChunks},
		{"twoTokens", []Token{"lorem", "ipsum"}, noChunks},
		{"threeTokens", []Token{"lorem", "ipsum", "dolor"}, oneChunk},
		{"multipleTokens", []Token{"lorem", "ipsum", "dolor", "sit", "amet"}, multiChunks},
	}
	for _, test := range cases {
		t.Run(test.name, func(tc *testing.T) {
			chunks := make([]Chunk, 0, len(test.want))
			for chunk := range getChunks(tokenReader(test.data)) {
				chunks = append(chunks, chunk)
			}
			assert.ElementsMatchf(tc, test.want, chunks, "Elements should match")
		})
	}
}

func TestCountChunks(t *testing.T) {
	chunk_a := Chunk{"a", "b", "c"}
	chunk_b := Chunk{"b", "c", "d"}
	cases := []struct {
		name string
		data []Chunk
		want ChunkFrequency
	}{
		{"simpleFreq", []Chunk{chunk_a, chunk_b, chunk_a}, ChunkFrequency{Count{2, chunk_a}, Count{1, chunk_b}}},
	}
	for _, test := range cases {
		t.Run(test.name, func(tc *testing.T) {
			freqs := countChunks(chunkReader(test.data))
			assert.ElementsMatch(tc, test.want, freqs, "Elements should match")
		})
	}
}

// I'm sure there's an easier way to push elements of a slice into a channel,
// but this works for testing
// TODO: Check if this can be made generic with interface{} and reflection
func tokenReader(data []Token) <-chan Token {
	ch := make(chan Token, ChunkSize)
	go func() {
		for _, tok := range data {
			ch <- tok
		}
		close(ch)
	}()
	return ch
}

func chunkReader(data []Chunk) <-chan Chunk {
	ch := make(chan Chunk)
	go func() {
		for _, chunk := range data {
			ch <- chunk
		}
		close(ch)
	}()
	return ch
}
