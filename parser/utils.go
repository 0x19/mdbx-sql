package parser

import (
	"strings"
	"unicode"
)

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_'
}

func lookupIdent(ident string) TokenType {
	switch strings.ToUpper(ident) {
	case "SELECT":
		return SELECT
	case "FROM":
		return FROM
	case "WHERE":
		return WHERE
	default:
		return IDENT
	}
}
