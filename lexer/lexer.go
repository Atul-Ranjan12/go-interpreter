package lexer

import (
	"strconv"

	"github.com/Atul-Ranjan12/errorHandler"
	"github.com/Atul-Ranjan12/token"
)

// Lexer defines the scanner to get tokens
type Lexer struct {
	Source string
	Tokens []*token.Token

	// Fields to keep track of where the Lexer
	// is on the code
	Start        int
	Current      int
	Line         int
	ErrorHandler errorHandler.ErrorHandler
}

// NewLexer returns an instance of scanner
func NewLexer(source string, errorHandler errorHandler.ErrorHandler) *Lexer {
	return &Lexer{
		Source:       source,
		Tokens:       make([]*token.Token, 0),
		Start:        0,
		Current:      0,
		Line:         1,
		ErrorHandler: errorHandler,
	}
}

// Advance function gets the character at Current
func (s *Lexer) Advance() byte {
	if s.IsAtEnd() {
		return 0
	}

	char := s.Source[s.Current]
	s.Current++
	return char
}

// IsAtEnd checks if the scanner has reached the end
func (s *Lexer) IsAtEnd() bool {
	return s.Current >= len(s.Source)
}

// AddToken adds a token
func (s *Lexer) AddToken(tokenType token.TokenType, literal interface{}) {
	text := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, &token.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    s.Line,
	})
}

// Match matches operators
func (s *Lexer) Match(expected byte) bool {
	if s.IsAtEnd() {
		return false
	}
	if s.Source[s.Current] != expected {
		return false
	}
	s.Current++
	return true
}

// Peek returns the next character
func (s *Lexer) Peek() byte {
	if s.IsAtEnd() {
		return 0
	}
	return s.Source[s.Current]
}

// IsDigit returns if a string is a digit
func (s *Lexer) IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsAlpha returns if string is an alphabet
func (s *Lexer) IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// IsAlphaNumeric returns if string is an alphanumeric
func (s *Lexer) IsAlphaNumeric(c byte) bool {
	return s.IsAlpha(c) || s.IsDigit(c)
}

// Get the next character in the source
func (s *Lexer) PeekNext() byte {
	if s.Current+1 >= len(s.Source) {
		return 0
	}
	return s.Source[s.Current+1]
}

// String function reads the entire string
func (s *Lexer) String() {
	// Within quotes and the file has not ended
	for s.Peek() != '"' && !s.IsAtEnd() {
		// Newline just simply increase the line
		if s.Peek() == '\n' {
			s.Line++
		}
		_ = s.Advance() // Keeps on increasing current
	}

	// Got an error, reached the end without closing the quote
	if s.IsAtEnd() {
		// TODO: Error handling
		// Lox.error(s.Line, "Unterminated string.")
		return
	}

	// The closing ".
	s.Advance()

	// Trim the surrounding quotes.
	value := s.Source[s.Start+1 : s.Current-1]
	s.AddToken(token.STRING, value)
}

// Number function reads the number
func (s *Lexer) Number() {
	for s.IsDigit(s.Peek()) {
		s.Advance()
	}

	// Look for a fractional part.
	if s.Peek() == '.' && s.IsDigit(s.PeekNext()) {
		// Consume the "."
		s.Advance()

		for s.IsDigit(s.Peek()) {
			s.Advance()
		}
	}

	value, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
	if err != nil {
		// TODO: Handle Error here
	}
	s.AddToken(token.NUMBER, value)
}

// Identifier checks if the token is an identifier
func (s *Lexer) Identifier() {
	for s.IsAlphaNumeric(s.Peek()) {
		s.Advance()
	}

	// Implement this function
	text := s.Source[s.Start:s.Current]
	keyword, ok := token.Keywords[text]
	if !ok {
		s.AddToken(token.IDENTIFIER, nil)
		return
	}

	s.AddToken(keyword, nil)
}

// ScanToken scans for a single token the tokens
func (s *Lexer) ScanToken() {
	c := s.Advance()
	switch c {
	case '(':
		s.AddToken(token.LEFT_PAREN, nil)
	case ')':
		s.AddToken(token.RIGHT_PAREN, nil)
	case '{':
		s.AddToken(token.LEFT_BRACE, nil)
	case '}':
		s.AddToken(token.RIGHT_BRACE, nil)
	case ',':
		s.AddToken(token.COMMA, nil)
	case '.':
		s.AddToken(token.DOT, nil)
	case '-':
		s.AddToken(token.MINUS, nil)
	case '+':
		s.AddToken(token.PLUS, nil)
	case ';':
		s.AddToken(token.SEMICOLON, nil)
	case '*':
		s.AddToken(token.STAR, nil)
	case '!':
		if s.Match('=') {
			s.AddToken(token.BANG_EQUAL, nil)
		} else {
			s.AddToken(token.BANG, nil)
		}
	case '=':
		if s.Match('=') {
			s.AddToken(token.EQUAL_EQUAL, nil)
		} else {
			s.AddToken(token.EQUAL, nil)
		}
	case '<':
		if s.Match('=') {
			s.AddToken(token.LESS_EQUAL, nil)
		} else {
			s.AddToken(token.LESS, nil)
		}
	case '>':
		if s.Match('=') {
			s.AddToken(token.GREATER_EQUAL, nil)
		} else {
			s.AddToken(token.GREATER, nil)
		}
	case '/':
		if s.Match('/') {
			// A comment goes until the end of the line.
			for s.Peek() != '\n' && !s.IsAtEnd() {
				// Just read the characters, do nothing
				_ = s.Advance()
			}
		} else {
			s.AddToken(token.SLASH, nil)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		s.Line++
	case '"':
		s.String()
	default:
		if s.IsDigit(c) {
			s.Number()
		} else if s.IsAlpha(c) {
			// Starts with an alphabet so
			// is an identifier
			s.Identifier()
		} else {
			// TODO: Error handling
			// Lox.error(s.Line, "Unexpected character.")
		}
	}
}

// TODO Implement ScanTokens
func (s *Lexer) ScanTokens() []*token.Token {
	for !s.IsAtEnd() {
		s.Start = s.Current
		s.ScanToken()
	}

	// Add the end of file to the tokens
	s.Tokens = append(s.Tokens, &token.Token{
		Type:    token.EOF,
		Literal: nil,
		Line:    s.Line,
	})

	return s.Tokens
}
