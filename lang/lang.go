package lang

import (
	"fmt"

	"github.com/Atul-Ranjan12/errorHandler"
	"github.com/Atul-Ranjan12/interpreter"
	"github.com/Atul-Ranjan12/lexer"
	"github.com/Atul-Ranjan12/parser"
)

type Lang struct {
	HadError    bool
	Lexer       *lexer.Lexer // The language has a lexer
	Parser      *parser.Parser
	Interpreter *interpreter.Interpreter
}

// NewLang initializes an instance of lang
func NewLang(source string) *Lang {
	lang := &Lang{}
	lang.Lexer = lexer.NewLexer(source, lang)

	// Lex the source for tokens
	tokens := lang.Lexer.ScanTokens()
	// Parse the tokens
	lang.Parser = parser.NewParser(tokens)

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
