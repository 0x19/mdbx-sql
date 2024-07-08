package mdbxsql

import (
	"errors"
	"fmt"
	"github.com/0x19/mdbx-sql/parser"
	"github.com/erigontech/mdbx-go/mdbx"
	"strings"
)

func Get(db *Db, sqlQuery string, record Model) error {
	lexer := parser.NewLexer(sqlQuery)
	sqlParser := parser.NewParser(lexer)
	stmt, ok := sqlParser.Parse().(*parser.SelectStatement)
	if !ok {
		return errors.New("invalid SQL query")
	}

	if stmt.TableName != "users" { // assuming we're only dealing with "users" table
		return errors.New("unsupported table")
	}

	fmt.Printf("AST: %+v\n", stmt)

	primaryKey, err := extractPrimaryKey(stmt.Conditions)
	if err != nil {
		return err
	}

	var value []byte
	err = db.env.View(func(txn *mdbx.Txn) error {
		value, err = txn.Get(db.dbi, primaryKey)
		return err
	})

	if err != nil {
		return err
	}

	return record.Unmarshal(value)
}

func extractPrimaryKey(conditions []parser.Condition) ([]byte, error) {
	for _, cond := range conditions {
		if strings.ToLower(cond.Field) == "id" && cond.Op == "=" {
			return []byte(cond.Value), nil
		}
	}
	return nil, errors.New("no valid primary key condition found")
}
