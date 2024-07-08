package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		// Basic SELECT statement
		{
			input: "SELECT name, age FROM users WHERE name = 'John' AND age = 30",
			expected: []Token{
				{Type: SELECT, Literal: "SELECT"},
				{Type: IDENT, Literal: "name"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
				{Type: FROM, Literal: "FROM"},
				{Type: IDENT, Literal: "users"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "John"},
				{Type: AND, Literal: "AND"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "="},
				{Type: NUMBER, Literal: "30"},
				{Type: EOF, Literal: ""},
			},
		},
		// INSERT statement
		{
			input: "INSERT INTO users (name, age) VALUES ('Jane', 25)",
			expected: []Token{
				{Type: INSERT, Literal: "INSERT"},
				{Type: INTO, Literal: "INTO"},
				{Type: IDENT, Literal: "users"},
				{Type: LPAREN, Literal: "("},
				{Type: IDENT, Literal: "name"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
				{Type: RPAREN, Literal: ")"},
				{Type: VALUES, Literal: "VALUES"},
				{Type: LPAREN, Literal: "("},
				{Type: IDENT, Literal: "Jane"},
				{Type: COMMA, Literal: ","},
				{Type: NUMBER, Literal: "25"},
				{Type: RPAREN, Literal: ")"},
				{Type: EOF, Literal: ""},
			},
		},
		// UPDATE statement
		{
			input: "UPDATE users SET name = 'Doe' WHERE id = 1",
			expected: []Token{
				{Type: UPDATE, Literal: "UPDATE"},
				{Type: IDENT, Literal: "users"},
				{Type: SET, Literal: "SET"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "Doe"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "id"},
				{Type: EQ, Literal: "="},
				{Type: NUMBER, Literal: "1"},
				{Type: EOF, Literal: ""},
			},
		},
		// DELETE statement
		{
			input: "DELETE FROM users WHERE id = 2",
			expected: []Token{
				{Type: DELETE, Literal: "DELETE"},
				{Type: FROM, Literal: "FROM"},
				{Type: IDENT, Literal: "users"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "id"},
				{Type: EQ, Literal: "="},
				{Type: NUMBER, Literal: "2"},
				{Type: EOF, Literal: ""},
			},
		},
		// Mixed whitespace
		{
			input: "SELECT\tname,\nage\nFROM\rusers WHERE\r\nname = 'John'\tAND age = 30",
			expected: []Token{
				{Type: SELECT, Literal: "SELECT"},
				{Type: IDENT, Literal: "name"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
				{Type: FROM, Literal: "FROM"},
				{Type: IDENT, Literal: "users"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "John"},
				{Type: AND, Literal: "AND"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "="},
				{Type: NUMBER, Literal: "30"},
				{Type: EOF, Literal: ""},
			},
		},
		// Handling semicolon
		{
			input: "SELECT * FROM users WHERE name = 'John' AND age = 30;",
			expected: []Token{
				{Type: SELECT, Literal: "SELECT"},
				{Type: ASTERISK, Literal: "*"},
				{Type: FROM, Literal: "FROM"},
				{Type: IDENT, Literal: "users"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "John"},
				{Type: AND, Literal: "AND"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "="},
				{Type: NUMBER, Literal: "30"},
				{Type: EOF, Literal: ";"},
			},
		},
		// JOIN statement
		{
			input: "SELECT users.name, orders.id FROM users JOIN orders ON users.id = orders.user_id WHERE users.name = 'John'",
			expected: []Token{
				{Type: SELECT, Literal: "SELECT"},
				{Type: IDENT, Literal: "users.name"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "orders.id"},
				{Type: FROM, Literal: "FROM"},
				{Type: IDENT, Literal: "users"},
				{Type: JOIN, Literal: "JOIN"},
				{Type: IDENT, Literal: "orders"},
				{Type: ON, Literal: "ON"},
				{Type: IDENT, Literal: "users.id"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "orders.user_id"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "users.name"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "John"},
				{Type: EOF, Literal: ""},
			},
		},
		// SUM aggregate function
		{
			input: "SELECT SUM(amount) FROM transactions WHERE status = 'complete';",
			expected: []Token{
				{Type: SELECT, Literal: "SELECT"},
				{Type: SUM, Literal: "SUM"},
				{Type: LPAREN, Literal: "("},
				{Type: IDENT, Literal: "amount"},
				{Type: RPAREN, Literal: ")"},
				{Type: FROM, Literal: "FROM"},
				{Type: IDENT, Literal: "transactions"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "status"},
				{Type: EQ, Literal: "="},
				{Type: IDENT, Literal: "complete"},
				{Type: EOF, Literal: ";"},
			},
		},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		for i, expectedToken := range tt.expected {
			tok := l.NextToken()

			require.Equal(
				t,
				expectedToken.Type,
				tok.Type,
				"tests[%d] - tokentype wrong. expected=%q, got=%q", i, expectedToken.Type, tok.Type,
			)
			require.Equal(
				t,
				expectedToken.Literal,
				tok.Literal,
				"tests[%d] - literal wrong. expected=%q, got=%q", i, expectedToken.Literal, tok.Literal,
			)

			// Debug prints
			//t.Logf("tests[%d] - got token: %+v", i, tok)
		}
	}
}
