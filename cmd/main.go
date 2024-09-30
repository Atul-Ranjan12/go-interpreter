package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Atul-Ranjan12/lang"
	"github.com/Atul-Ranjan12/parser/astprinter"
	"github.com/Atul-Ranjan12/token"
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
	expression, err := lang.Parser.Parse()
	if err != nil {
		log.Println("An error occured while parsing: ", err)
		return
	}

	// Print each token
	for _, t := range lang.Lexer.Tokens {
		fmt.Printf("%s %s %v\n", token.TokenTypeToString(t.Type), t.Lexeme, t.Literal)
	}

	log.Println("\n\nPrinting the output from the parser")

	printer := astprinter.NewAstPrinter()
	fmt.Println(printer.Print(expression))

	// Evaluating the output
	log.Println("\n\nInterpreting the expressions")
	result, err := lang.Interpreter.Interpret(expression)
	if err != nil {
		log.Println("Error in interpretation: ", err)
		return
	}

	log.Println(result)
}
