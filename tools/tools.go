package tools

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// This package contains all the
// tools for the interpreter

// Run Function runs the file
func Run(content []byte) error {
	// 1. Initialize the scanner
	// 2. Do Stuff
	//
	log.Println(content)
	return nil
}

// Function to run the file
func RunFile(path string) error {
	// Read the file from the path
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading the file contents")
		return err
	}

	Run(content)

	return nil
}

// Function to run the prompt
func RunPrompt() error {
	fmt.Println("Starting execution")
	// Initialize the reader
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("An error occured reading the prompt, try again: ", err)
			continue
		}

		if input == "quit()" {
			fmt.Println("Exiting...")
			break
		}

		// Process the input
		Run([]byte(input))
	}

	return nil
}
