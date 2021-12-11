package main

import (
	"bufio"
	"fmt"
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

	// Scan file lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Notify if some error happened during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return 2
	}

	return 0
}
