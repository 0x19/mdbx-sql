package mdbxsql

import (
	"errors"
	"github.com/erigontech/mdbx-go/mdbx"
)

type Table struct {
	Name    string
	Primary string
	db      *Db
	dbi     mdbx.DBI
}

type Schema struct {
	tables map[string]*Table
	db     *Db
}

func NewSchema(db *Db) *Schema {
	return &Schema{
		tables: make(map[string]*Table),
		db:     db,
	}
}

func (s *Schema) CreateTable(name string, primary string) (*Table, error) {
	if name == "" {
		return nil, errors.New("table name cannot be empty")
	}

	var dbi mdbx.DBI
	err := s.db.env.Update(func(txn *mdbx.Txn) (err error) {
		dbi, err = txn.OpenDBI(name, mdbx.Create, nil, nil)
		return err
	})

	if err != nil {
		return nil, err
	}

	table := &Table{
		Name:    name,
		Primary: primary,
		db:      s.db,
		dbi:     dbi,
	}
	s.tables[name] = table
	return table, nil
}

func (s *Schema) GetTable(name string) (*Table, error) {
	table, exists := s.tables[name]
	if !exists {
		return nil, errors.New("table not found")
	}
	return table, nil
}
