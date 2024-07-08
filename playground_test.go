package mdbxsql

import (
	"bytes"
	"fmt"
	"github.com/0x19/mdbx-sql/parser"
	"github.com/golang/snappy"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"testing"
	"time"
)

type UserGo struct {
	ID   int32  `msgpack:"id"`
	Name string `msgpack:"name"`
	Age  int32  `msgpack:"age"`
}

func (u *UserGo) PrimaryKey() []byte {
	return []byte(strconv.FormatInt(int64(u.ID), 10))
}

func (u *UserGo) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	err := encoder.Encode(u)
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

	decoder := msgpack.NewDecoder(bytes.NewBuffer(decompressed))
	err = decoder.Decode(u)
	if err != nil {
		return err
	}

	return nil
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
