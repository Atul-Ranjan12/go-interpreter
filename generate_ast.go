package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type AstType struct {
	name   string
	fields []string
}

func defineAst(outputDir, baseName string, types []AstType) error {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	write := func(format string, args ...interface{}) {
		fmt.Fprintf(file, format, args...)
	}

	write("package expressions\n\n")
	write("import \"github.com/Atul-Ranjan12/token\"\n\n")

	// Define the base interface (Expr or Stmt)
	write("type %s interface {\n", baseName)
	write("\tAccept(visitor %sVisitor) (interface{}, error)\n", baseName)
	write("}\n\n")

	// Define the Visitor interface
	defineVisitor(write, baseName, types)

	// Define each AST struct
	for _, t := range types {
		defineType(write, baseName, t)
	}

	return nil
}

func defineVisitor(write func(string, ...interface{}), baseName string, types []AstType) {
	write("type %sVisitor interface {\n", baseName)
	for _, t := range types {
		typeName := t.name
		write("\tVisit%s%s(%s *%s) (interface{}, error)\n", typeName, baseName, strings.ToLower(baseName), typeName)
	}
	write("}\n\n")
}

func defineType(write func(string, ...interface{}), baseName string, t AstType) {
	write("// These are functions for %s \n", t.name)

	write("type %s struct {\n", t.name)
	for _, field := range t.fields {
		write("\t%s\n", field)
	}
	write("}\n\n")

	write("var _ %s = (*%s)(nil)\n\n", baseName, t.name)

	// Define Accept method
	write("func (e *%s) Accept(visitor %sVisitor) (interface{}, error) {\n", t.name, baseName)
	write("\treturn visitor.Visit%s%s(e)\n", t.name, baseName)
	write("}\n\n")
}

func main() {
	outputDir := "parser/expressions"

	err := defineAst(outputDir, "Expr", []AstType{
		{"Assign", []string{"Name token.Token", "Value Expr"}},
		{"Logical", []string{"Left Expr", "Right Expr", "Operator token.Token"}},
		{"Binary", []string{"Left Expr", "Operator token.Token", "Right Expr"}},
		{"Grouping", []string{"Expression Expr"}},
		{"Literal", []string{"Value interface{}"}},
		{"Unary", []string{"Operator token.Token", "Right Expr"}},
		{"Variable", []string{"Name token.Token"}},
		{"BreakExpr", []string{}},
	})
	if err != nil {
		log.Fatalf("Error generating Expr AST: %v", err)
	}
	err = defineAst(outputDir, "Stmt", []AstType{
		{"Block", []string{"Statements []Stmt"}},
		{"ExprStatement", []string{"Expression Expr"}},
		{"PrintStatement", []string{"Expression Expr"}},
		{"WhileStatement", []string{"Condition Expr", "Body Stmt"}},
		{"Var", []string{"Name token.Token", "Initializer Expr"}},
		{"If", []string{"Condition Expr", "ThenBranch Stmt", "ElseBranch Stmt"}},
	})
	if err != nil {
		log.Fatalf("Error generating Expr AST: %v", err)
	}

	log.Println("Successfully generated ASTs at:", outputDir)
}
