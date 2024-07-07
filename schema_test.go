package mdbxsql

import (
	"context"
	"fmt"
	"github.com/0x19/mdbx-sql/parser"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, u)
}

func TestParserAndDatabase(t *testing.T) {
	// Test SQL Parsing
	input := "SELECT name, age FROM users WHERE active"
	start := time.Now()
	lexer := parser.NewLexer(input)
	ddl := parser.NewParser(lexer)

	ast, err := ddl.Parse()
	require.NoError(t, err)
	log.Printf("SQL Parsing completed in %v", time.Since(start))
	log.Printf("AST: %+v", ast)

	// Test MDBX Database Operations
	ctx := context.Background()
	db, err := NewDb(ctx, "testdb")
	require.NoError(t, err)
	defer db.Close()

	schema := NewSchema(db)
	userTable := schema.CreateTable("users", "ID")

	// Test Insert
	user := &User{ID: 1, Name: "John Doe", Age: 30}
	start = time.Now()
	err = Insert(userTable, user)
	require.NoError(t, err)
	log.Printf("Insert operation completed in %v", time.Since(start))

	// Test Get
	start = time.Now()
	retrievedUser := &User{}
	err = Get(userTable, 1, retrievedUser)
	require.NoError(t, err)
	require.Equal(t, user, retrievedUser)
	log.Printf("Get operation completed in %v", time.Since(start))
	log.Printf("Retrieved User: %+v", retrievedUser)

	// Test Update
	user.Age = 31
	start = time.Now()
	err = Update(userTable, user)
	require.NoError(t, err)
	log.Printf("Update operation completed in %v", time.Since(start))

	start = time.Now()
	retrievedUser = &User{}
	err = Get(userTable, 1, retrievedUser)
	require.NoError(t, err)
	require.Equal(t, 31, retrievedUser.Age)
	log.Printf("Get operation (post-update) completed in %v", time.Since(start))
	log.Printf("Updated User: %+v", retrievedUser)

	// Test Delete
	start = time.Now()
	err = Delete(userTable, 1)
	require.NoError(t, err)
	log.Printf("Delete operation completed in %v", time.Since(start))

	start = time.Now()
	retrievedUser = &User{}
	err = Get(userTable, 1, retrievedUser)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
	log.Printf("Get operation (post-delete) completed in %v", time.Since(start))
}

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
