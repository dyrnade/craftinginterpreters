package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var hadError = false

func main() {
	args := os.Args

	if len(args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func runFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer file.Close()

	// Create a bufio.Reader from the file
	reader := bufio.NewReader(file)

	// Create a new CustomScanner instance with the input source
	inputString := readStringFromReader(reader)
	scanner := NewCustomScanner(inputString)

	// Scan tokens
	tokens := scanner.ScanTokens()

	// Print the scanned tokens
	for _, token := range tokens {
		fmt.Println(token.String())
	}

	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	input := os.Stdin
	reader := bufio.NewReader(input)

	for {
		fmt.Print("> ")
		inputString := readStringFromReader(reader)
		scanner := NewCustomScanner(inputString)

		tokens := scanner.ScanTokens()

		for _, token := range tokens {
			fmt.Println(token)
		}

		// run()
		hadError = false
	}
}

func readStringFromReader(reader *bufio.Reader) string {
	var sb strings.Builder
	for {
		line, err := reader.ReadString('\n')
		sb.WriteString(line)
		if err != nil {
			break
		}
	}
	return sb.String()
}

// func run(s string) {
// 	nwc := *NewCustomScanner(s)
// 	nwc.ScanTokens()
// }
