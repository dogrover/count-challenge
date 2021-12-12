package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

	// Scan file for words. Word delimiters are defined by unicode.isSpace.
	// scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanWords)
	// for scanner.Scan() {
	// 	word := strings.TrimFunc(scanner.Text(), func(r rune) bool {
	// 		return !unicode.IsLetter(r)
	// 	})
	// 	if len(word) > 0 {
	// 		fmt.Println(word)
	// 	}
	// }

	// // Notify if some error happened during scanning
	// if err := scanner.Err(); err != nil {
	// 	fmt.Println(err)
	// 	return 2
	// }

	for word := range readTokens(file) {
		fmt.Println(word)
	}

	return 0
}

type Token string

func readTokens(file io.Reader) chan Token {
	ch := make(chan Token)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	go func() {
		for scanner.Scan() {
			ch <- Token(scanner.Text())
		}
		close(ch)
	}()
	return ch
}
