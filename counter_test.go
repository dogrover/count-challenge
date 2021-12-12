package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadTokens(t *testing.T) {
	data := strings.NewReader("Lorem ipsum sit amet")
	want := []Token{"Lorem", "ipsum", "sit", "amet"}

	words := make([]Token, 0, 5)
	for w := range readTokens(data) {
		words = append(words, w)
	}
	assert.ElementsMatch(t, want, words, "Elements should match")
}
