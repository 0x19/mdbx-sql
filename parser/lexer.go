package parser

import (
	"unicode"
)

// Lexer represents a lexer for SQL.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// NewLexer initializes a new Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	//log.Printf("readChar: position=%d readPosition=%d ch=%q", l.position, l.readPosition, l.ch)
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	//log.Printf("NextToken: ch=%q", l.ch)

	switch l.ch {
	case ',':
		tok = newToken(COMMA, l.ch)
	case '=':
		tok = newToken(EQ, l.ch)
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(GT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(LT, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case ';':
		tok = newToken(EOF, l.ch)
		l.readChar()
		return tok
	case '\'':
		tok.Type = IDENT
		tok.Literal = l.readString()
		//log.Printf("NextToken: string literal=%s", tok.Literal)
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			//log.Printf("NextToken: identifier=%s type=%d", tok.Literal, tok.Type)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			//log.Printf("NextToken: number=%s", tok.Literal)
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1 // skip the opening single quote
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	stringValue := l.input[position:l.position]
	l.readChar() // skip the closing single quote
	return stringValue
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}
