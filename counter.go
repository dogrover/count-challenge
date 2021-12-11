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
	// CHeck that we've been given at least one file to process
	numArgs := len(os.Args[1:])
	if numArgs < 1 {
		fmt.Fprintln(os.Stderr, "No filenames given")
		return 1
	} else if numArgs > 1 {
		fmt.Fprintln(os.Stderr, "Multiple filenames given")
		return 1
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return 0
}
