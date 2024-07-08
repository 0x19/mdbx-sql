package parser

import "sync"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	WS
	IDENT
	NUMBER
	COMMA
	SELECT
	FROM
	WHERE
	EQ
	AND
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"SELECT": SELECT,
	"FROM":   FROM,
	"WHERE":  WHERE,
	"AND":    AND,
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
