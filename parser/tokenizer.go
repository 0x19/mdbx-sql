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
	GT  // Greater Than
	LT  // Less Than
	GTE // Greater Than or Equal To
	LTE // Less Than or Equal To
	NEQ // Not Equal To
	AND
	INSERT
	INTO
	VALUES
	UPDATE
	SET
	DELETE
	LPAREN
	RPAREN
	JOIN
	ON
	SUM
	ASTERISK
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
	"INSERT": INSERT,
	"INTO":   INTO,
	"VALUES": VALUES,
	"UPDATE": UPDATE,
	"SET":    SET,
	"DELETE": DELETE,
	"JOIN":   JOIN,
	"ON":     ON,
	"SUM":    SUM,
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
