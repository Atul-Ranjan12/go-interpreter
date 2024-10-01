package main

import (
	"io/ioutil"
	"log"

	"github.com/Atul-Ranjan12/lang"
)

// Interperter Main

// func main() {
// 	// os.Args[0] is the program name
// 	// everything after that contains the arguments
// 	// Has more than the file to interpret
// 	if len(os.Args) > 2 {
// 		log.Fatal("Unexpected arguments")
// 		return
// 	} else if len(os.Args) == 2 {
// 		// Here we run the file
// 		tools.RunFile(os.Args[1])
// 	} else {
// 		// Here we run the prompt
// 		tools.RunPrompt()
// 	}
// }

func main() {
	// Read the contents of test.rnj
	content, err := ioutil.ReadFile("test.rnj")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Create a new lexer with the file contents
	lang := lang.NewLang(string(content))

	// Parse
	statements, err := lang.Parser.Parse()
	if err != nil {
		log.Println("An error occured while parsing: ", err)
		return
	}

	log.Println("Printing statements: ", statements)
	log.Println("Successful parse ")
	log.Println("\nOutput")

	err = lang.Interpreter.Interpret(statements)
	if err != nil {
		log.Println("Interpretation Error: ", err)
	}

	log.Println("Program execution successful")

}
