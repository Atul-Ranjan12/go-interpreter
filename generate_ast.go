package main

import (
	"fmt"
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

	// Define the Expr interface
	write("type Expr interface {\n")
	write("\tAccept(visitor ExprVisitor) (interface{}, error)\n")
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
	write("type ExprVisitor interface {\n")
	for _, t := range types {
		typeName := t.name
		write("\tVisit%s%s(%s *%s) (interface{}, error)\n", typeName, baseName, strings.ToLower(baseName), typeName)
	}
	write("}\n\n")
}

func defineType(write func(string, ...interface{}), baseName string, t AstType) {
	write("// Tese are functions for %s \n", t.name)

	write("type %s struct {\n", t.name)
	for _, field := range t.fields {
		write("\t%s\n", field)
	}
	write("}\n\n")

	write("var _ Expr = (*%s)(nil)\n\n", t.name)

	// Define Accept method
	write("func (e *%s) Accept(visitor ExprVisitor) (interface{}, error) {\n", t.name)
	write("\treturn visitor.Visit%s%s(e)\n", t.name, baseName)
	write("}\n\n")
}

func main() {
	outputDir := "parser/expressions"
	defineAst(outputDir, "Expr", []AstType{
		{"Binary", []string{"Left Expr", "Operator token.Token", "Right Expr"}},
		{"Grouping", []string{"Expression Expr"}},
		{"Literal", []string{"Value interface{}"}},
		{"Unary", []string{"Operator token.Token", "Right Expr"}},
	})
}
