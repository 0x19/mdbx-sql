package mdbxsql

import (
	"bytes"
	"capnproto.org/go/capnp/v3"
	"context"
	"fmt"
	"github.com/0x19/mdbx-sql/parser"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

type UserGo struct {
	ID   int32
	Name string
	Age  int32
}

func (u *UserGo) PrimaryKey() interface{} {
	return u.ID
}

func (u *UserGo) Marshal() ([]byte, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	userCapnp, err := NewRootUser(seg)
	if err != nil {
		return nil, err
	}

	userCapnp.SetId(u.ID)
	userCapnp.SetName(u.Name)
	userCapnp.SetAge(u.Age)

	buf := new(bytes.Buffer)
	err = capnp.NewEncoder(buf).Encode(msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *UserGo) Unmarshal(data []byte) error {
	msg, err := capnp.NewDecoder(bytes.NewBuffer(data)).Decode()
	if err != nil {
		return err
	}

	userCapnp, err := ReadRootUser(msg)
	if err != nil {
		return err
	}

	u.ID = userCapnp.Id()
	u.Name, err = userCapnp.Name()
	if err != nil {
		return err
	}
	u.Age = userCapnp.Age()

	return nil
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
	user := &UserGo{ID: 1, Name: "John Doe", Age: 30}
	start = time.Now()
	err = Insert(userTable, user)
	require.NoError(t, err)
	log.Printf("Insert operation completed in %v", time.Since(start))

	// Test Get
	start = time.Now()
	retrievedUser := &UserGo{}
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
	retrievedUser = &UserGo{}
	err = Get(userTable, 1, retrievedUser)
	require.NoError(t, err)
	require.Equal(t, int32(31), retrievedUser.Age)
	log.Printf("Get operation (post-update) completed in %v", time.Since(start))
	log.Printf("Updated User: %+v", retrievedUser)

	// Test Delete
	start = time.Now()
	err = Delete(userTable, 1)
	require.NoError(t, err)
	log.Printf("Delete operation completed in %v", time.Since(start))

	start = time.Now()
	retrievedUser = &UserGo{}
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
