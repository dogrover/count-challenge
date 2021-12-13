package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func main() {
	// Exit halts execution before deferred functions, log flushes, and other
	// cleanup. So keep those in a separate function.
	os.Exit(run())
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
	for tok := range tokens {
		fmt.Println(tok)
	}

	return 0
}

type Word string
type Token string

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
