package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadTokens(t *testing.T) {
	expectedTokens := []Token{"Lorem", "ipsum", "sit", "amet"}
	cases := []struct {
		name string
		data string
		want []Token
	}{
		{"simpleTokens", "Lorem ipsum sit amet", expectedTokens},
		{"withNewline", "Lorem\nipsum\nsit\namet", expectedTokens},
		{"extraSpace", "    Lorem  \n  ipsum\n\n\nsit amet     ", expectedTokens},
	}
	for _, test := range cases {
		data := strings.NewReader(test.data)
		words := make([]Token, 0, len(test.want))
		for w := range readTokens(data) {
			words = append(words, w)
		}
		assert.ElementsMatch(t, test.want, words, "%v: Elements should match", test.name)
	}
}
