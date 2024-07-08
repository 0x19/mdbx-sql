package parser

import "sync"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	WS
	IDENT
	COMMA
	SELECT
	FROM
	WHERE
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
}

var tokenPool = sync.Pool{
	New: func() interface{} {
		return &Token{}
	},
}

func (tok *Token) reset() {
	tok.Type = ILLEGAL
	tok.Literal = ""
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
