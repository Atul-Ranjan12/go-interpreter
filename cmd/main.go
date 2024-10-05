package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Atul-Ranjan12/lang"
	"github.com/Atul-Ranjan12/parser/astprinter"
	"github.com/Atul-Ranjan12/parser/expressions"
)

func PrintAST(stmt expressions.Stmt, depth int) {
	printer := astprinter.NewAstPrinter()
	result, err := stmt.Accept(printer)
	if err != nil {
		log.Printf("Error printing AST: %v", err)
		return
	}

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s\n", indent, result)

	switch s := stmt.(type) {
	case *expressions.Block:
		for _, subStmt := range s.Statements {
			PrintAST(subStmt, depth+1)
		}
	case *expressions.If:
		PrintAST(s.ThenBranch, depth+1)
		if s.ElseBranch != nil {
			PrintAST(s.ElseBranch, depth+1)
		}
	case *expressions.WhileStatement:
		PrintAST(s.Body, depth+1)
	}
}

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
	content, err := ioutil.ReadFile("test.lang")
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

	fmt.Println("AST Structure:")
	for _, statement := range statements {
		PrintAST(statement, 0)
	}

	// log.Println("Printing statements: ", statements)
	// log.Println("Successful parse ")
	// log.Println("\nOutput")

	// Resolve it here
	err = lang.Resolver.ResolveStatements(statements)
	if err != nil {
		log.Println("Error in resolving the code: ", err)
		return
	}

	err = lang.Interpreter.Interpret(statements)
	if err != nil {
		log.Println("Interpretation Error: ", err)
	}

	log.Println("Program execution successful")

}
