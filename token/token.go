package token

import (
	"fmt"
)

// These are the standard Keywords for the language
var Keywords = map[string]TokenType{
	"and":     AND,
	"break":   BREAK,
	"struct":  CLASS,
	"else":    ELSE,
	"false":   FALSE,
	"for":     FOR,
	"def":     FUN,
	"if":      IF,
	"nil":     NIL,
	"or":      OR,
	"println": PRINT,
	"return":  RETURN,
	"super":   SUPER,
	"this":    THIS,
	"true":    TRUE,
	"var":     VAR,
	"while":   WHILE,
}

// Token represents a token in the source code
type Token struct {
	Type    TokenType // This is the type of token
	Lexeme  string    //
	Literal interface{}
	Line    int
}

// NewToken creates a new Token instance
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	switch t.Type {
	case NUMBER, STRING:
		return fmt.Sprintf("%s %s %v", TokenTypeToString(t.Type), t.Lexeme, t.Literal)
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("%s %s", TokenTypeToString(t.Type), t.Lexeme)
	}
}
