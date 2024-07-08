package mdbxsql

import (
	"bytes"
	"capnproto.org/go/capnp/v3"
	"context"
	"fmt"
	"github.com/0x19/mdbx-sql/parser"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/require"
	"log"
	"strconv"
	"testing"
	"time"
)

type UserGo struct {
	ID   int32
	Name string
	Age  int32
}

func (u *UserGo) PrimaryKey() []byte {
	return []byte(strconv.FormatInt(int64(u.ID), 10))
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

	if err := userCapnp.SetName(u.Name); err != nil {
		return nil, err
	}

	userCapnp.SetAge(u.Age)

	buf := new(bytes.Buffer)
	err = capnp.NewEncoder(buf).Encode(msg)
	if err != nil {
		return nil, err
	}

	compressed := snappy.Encode(nil, buf.Bytes())
	return compressed, nil
}

func (u *UserGo) Unmarshal(data []byte) error {
	decompressed, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	msg, err := capnp.NewDecoder(bytes.NewBuffer(decompressed)).Decode()
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

	ast := ddl.Parse()
	log.Printf("SQL Parsing completed in %v", time.Since(start))
	log.Printf("AST: %+v", ast)

	// Test MDBX Database Operations
	ctx := context.Background()
	db, err := NewDb(ctx, "testdb", 10)
	require.NoError(t, err)
	defer db.Close()

	schema := NewSchema(db)
	userTable, err := schema.CreateTable("users", "ID")
	require.NoError(t, err)

	// Test Insert
	user := &UserGo{ID: 1, Name: "John Doe", Age: 30}
	start = time.Now()
	err = Insert(userTable, user)
	require.NoError(t, err)
	log.Printf("Insert operation completed in %v", time.Since(start))

	// Test Get
	start = time.Now()
	retrievedUser := &UserGo{}
	err = Get(userTable, []byte(strconv.FormatInt(1, 10)), retrievedUser)
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
	err = Get(userTable, []byte(strconv.FormatInt(1, 10)), retrievedUser)
	require.NoError(t, err)
	require.Equal(t, int32(31), retrievedUser.Age)
	log.Printf("Get operation (post-update) completed in %v", time.Since(start))
	log.Printf("Updated User: %+v", retrievedUser)

	// Test Delete
	start = time.Now()
	err = Delete(userTable, []byte(strconv.FormatInt(1, 10)))
	require.NoError(t, err)
	log.Printf("Delete operation completed in %v", time.Since(start))

	start = time.Now()
	retrievedUser = &UserGo{}
	err = Get(userTable, []byte(strconv.FormatInt(1, 10)), retrievedUser)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
	log.Printf("Get operation (post-delete) completed in %v", time.Since(start))
}

func TestPlayground(t *testing.T) {
	input := "SELECT name, age FROM users WHERE active"
	start := time.Now()
	lexer := parser.NewLexer(input)
	ddl := parser.NewParser(lexer)

	ast := ddl.Parse()
	fmt.Printf("AST: %+v in %v \n", ast, time.Since(start))
}

func BenchmarkPlayground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		input := "SELECT name, age FROM users WHERE active"
		lexer := parser.NewLexer(input)
		ddl := parser.NewParser(lexer)
		b.StartTimer()

		ast := ddl.Parse()
		_ = ast // discard the result to focus on the performance
	}
}
