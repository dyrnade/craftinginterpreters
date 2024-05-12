package main

import (
	"fmt"
	"strconv"
)

type CustomScanner struct {
	Source  string
	Start   int
	Current int
	Line    int
	Tokens  []Token
}

func NewCustomScanner(source string) *CustomScanner {
	return &CustomScanner{
		Source:  source,
		Start:   0,
		Current: 0,
		Line:    1,
		Tokens:  make([]Token, 0),
	}
}

func (cs *CustomScanner) ScanTokens() []Token {
	for !cs.isAtEnd() {
		cs.Start = cs.Current
		cs.scanToken()
	}

	cs.Tokens = append(cs.Tokens, *NewToken(EOF, "", nil, cs.Line))
	return cs.Tokens
}

func (cs *CustomScanner) isAtEnd() bool {
	return cs.Current >= len(cs.Source)
}

func (cs *CustomScanner) scanToken() {
	c := cs.advance()

	switch c {
	case '(':
		cs.addToken(LEFT_PAREN, nil)
	case ')':
		cs.addToken(RIGHT_PAREN, nil)
	case '{':
		cs.addToken(LEFT_BRACE, nil)
	case '}':
		cs.addToken(RIGHT_BRACE, nil)
	case ',':
		cs.addToken(COMMA, nil)
	case '.':
		cs.addToken(DOT, nil)
	case '-':
		cs.addToken(MINUS, nil)
	case '+':
		cs.addToken(PLUS, nil)
	case ';':
		cs.addToken(SEMICOLON, nil)
	case '*':
		cs.addToken(STAR, nil)
	case '!':
		if cs.match('=') {
			cs.addToken(BANG_EQUAL, nil)
		} else {
			cs.addToken(BANG, nil)
		}
	case '=':
		if cs.match('=') {
			cs.addToken(EQUAL_EQUAL, nil)
		} else {
			cs.addToken(EQUAL, nil)
		}
	case '<':
		if cs.match('=') {
			cs.addToken(LESS_EQUAL, nil)
		} else {
			cs.addToken(LESS, nil)
		}
	case '>':
		if cs.match('=') {
			cs.addToken(GREATER_EQUAL, nil)
		} else {
			cs.addToken(GREATER, nil)
		}
	case '/':
		if cs.match('/') {
			for cs.peek() != '\n' && !cs.isAtEnd() {
				cs.advance()
			}
		} else {
			cs.addToken(SLASH, nil)
		}
	case '"':
		cs.string()
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		cs.Line++
	default:
		if isDigit(c) {
			cs.number()
		} else if isAlpha(c) {
			cs.identifier()
		} else {
			cs.error(cs.Line, fmt.Sprintf("Unexpected character %s.\n", cs.Source[cs.Start:cs.Current]))
		}
	}
}

// 1234
// 12.34

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (cs *CustomScanner) identifier() {
	for isAlpha(cs.peek()) {
		cs.advance()
	}

	keywords := make(map[string]TokenType)

	// Add keywords and their TokenType values to the map
	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE

	if tokenType, ok := keywords[cs.Source[cs.Start:cs.Current]]; ok {
		cs.addToken(tokenType, nil)
	} else {
		cs.addToken(IDENTIFIER, nil)
	}
}

func (cs *CustomScanner) number() {
	if isDigit(cs.peek()) {
		for isDigit(cs.peek()) {
			cs.advance()
		}
	}

	if cs.peek() == '.' {
		cs.advance()
		for isDigit(cs.peek()) {
			cs.advance()
		}
	}

	n, err := strconv.ParseFloat(cs.Source[cs.Start:cs.Current], 64)
	if err != nil {
		cs.error(cs.Line, fmt.Sprintf("Can not parse to number %s\n", cs.Source[cs.Start:cs.Current]))
	}
	cs.addToken(NUMBER, n)
}

func (cs *CustomScanner) string() {
	for cs.peek() != '"' && !cs.isAtEnd() {
		if cs.peek() == '\n' {
			cs.Line++
		}
		cs.advance()
	}

	if cs.isAtEnd() {
		cs.error(cs.Line, "Unterminated string.")
		return
	}

	cs.advance() // The closing ".

	s := cs.Source[cs.Start+1 : cs.Current-1]
	cs.addToken(STRING, s)
}

func (cs *CustomScanner) peek() byte {
	if cs.isAtEnd() {
		return '\x00'
	}

	return cs.Source[cs.Current]
}

func (cs *CustomScanner) advance() byte {

	if cs.isAtEnd() {
		return '\x00'
	}
	c := cs.Source[cs.Current]
	cs.Current++
	return c
}

func (cs *CustomScanner) match(expected byte) bool {

	if cs.isAtEnd() {
		return false
	}

	if cs.Source[cs.Current] != expected {
		return false
	}

	cs.Current++
	return true
}

func (cs *CustomScanner) error(line int, message string) {
	cs.report(line, "", message)
}

func (cs *CustomScanner) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	hadError = true
}

func (cs *CustomScanner) addToken(tokenType TokenType, literal interface{}) {
	text := cs.Source[cs.Start:cs.Current]
	// fmt.Println("START: ", cs.Start)
	// fmt.Println("CURRENT: ", cs.Current)
	cs.Tokens = append(cs.Tokens, *NewToken(tokenType, text, literal, cs.Line))
}
