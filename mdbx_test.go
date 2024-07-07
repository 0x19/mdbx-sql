package mdbxsql

import (
	"fmt"
	"github.com/0x19/mdbx-sql/parser"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPlayground(t *testing.T) {
	input := "SELECT name, age FROM users WHERE active"
	start := time.Now()
	lexer := parser.NewLexer(input)
	ddl := parser.NewParser(lexer)

	ast, err := ddl.Parse()
	require.NoError(t, err)
	fmt.Printf("AST: %+v in %v \n", ast, time.Since(start))
}

func BenchmarkPlayground(b *testing.B) {
	input := "SELECT name, age FROM users WHERE active"
	lexer := parser.NewLexer(input)
	ddl := parser.NewParser(lexer)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		lexer.Init(input)
		ddl.Init(lexer)
		b.StartTimer()

		ast, _ := ddl.Parse()
		_ = ast // discard the result to focus on the performance
	}
}
