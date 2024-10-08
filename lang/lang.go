package lang

import (
	"fmt"

	"github.com/Atul-Ranjan12/errorHandler"
	"github.com/Atul-Ranjan12/interpreter"
	"github.com/Atul-Ranjan12/lexer"
	"github.com/Atul-Ranjan12/parser"
	"github.com/Atul-Ranjan12/resolver"
)

type Lang struct {
	HadError    bool
	Lexer       *lexer.Lexer // The language has a lexer
	Parser      *parser.Parser
	Resolver    *resolver.Resolver
	Interpreter *interpreter.Interpreter
}

// NewLang initializes an instance of lang
func NewLang(source string) *Lang {
	lang := &Lang{}
	lang.Lexer = lexer.NewLexer(source, lang)

	// Lex the source for tokens
	tokens := lang.Lexer.ScanTokens()
	// Initialize the parser
	lang.Parser = parser.NewParser(tokens)
	// Initialize the interpreter
	lang.Interpreter = interpreter.NewInterpreter()
	// Initialize the resolver
	lang.Resolver = resolver.NewResolver(lang.Interpreter)
	return lang
}

func (l *Lang) Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func (l *Lang) Error(line int, message string) {
	l.Report(line, "", message)
}

func (l *Lang) HasError() bool {
	return l.HadError
}

func (l *Lang) ResetError() {
	l.HadError = false
}

// Ensure Lang implements ErrorHandler
var _ errorHandler.ErrorHandler = (*Lang)(nil)
