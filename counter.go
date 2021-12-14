package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const ChunkSize = 3 // Number of tokens to index on

type Word string            // White-space-delimited unit of text
type Token string           // A lower-case word, trimmed of punctuation
type Chunk [ChunkSize]Token // A group of tokens

// Scans a reader for words, delimited by whitespace (as defined by unicode.isSpace)
func readWords(file io.Reader) <-chan Word {
	ch := make(chan Word)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	go func() {
		for scanner.Scan() {
			ch <- Word(scanner.Text())
		}
		close(ch)
	}()
	return ch
}

// Strips all leading and trailing non-letter characters from a word, and
// converts it to a lower case token.
func wordsToTokens(words <-chan Word) <-chan Token {
	ch := make(chan Token)
	go func() {
		for word := range words {
			tok := strings.TrimFunc(string(word), func(r rune) bool {
				return !unicode.IsLetter(r)
			})
			if len(tok) > 0 {
				ch <- Token(strings.ToLower(tok))
			}
		}
		close(ch)
	}()
	return ch
}

func getChunks(tokens <-chan Token) <-chan Chunk {
	ch := make(chan Chunk)
	var chunk Chunk

	// Fill the first chunk with tokens, if we have enough
	for i := 0; i < ChunkSize; i++ {
		tok, ok := <-tokens
		if !ok {
			close(ch)
			break
		}
		chunk[i] = tok
	}

	// If we can't fill the first chunk, we're done
	if chunk[ChunkSize-1] == "" {
		return ch
	}

	// Generate new chunks by popping the first token out of the chunk, and
	// pushing a new one onto the end
	go func() {
		for tok := range tokens {
			ch <- chunk
			copy(chunk[:], chunk[1:])
			chunk[ChunkSize-1] = tok
		}
		ch <- chunk
		close(ch)
	}()
	return ch
}

func run() int {
	// Check that we've been given exactly one file to process
	numArgs := len(os.Args[1:])
	if numArgs < 1 {
		fmt.Fprintln(os.Stderr, "No filenames given")
		return 1
	} else if numArgs > 1 {
		fmt.Fprintln(os.Stderr, "Multiple filenames given")
		return 1
	}

	// Check for valid file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return 2
	}
	defer file.Close()

	words := readWords(file)
	tokens := wordsToTokens(words)
	for chunk := range getChunks(tokens) {
		fmt.Println(chunk)
	}

	return 0
}

func main() {
	// Exit halts execution before deferred functions, log flushes, and other
	// cleanup. So keep those in a separate function.
	os.Exit(run())
}
