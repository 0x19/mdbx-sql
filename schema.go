package mdbxsql

import (
	"errors"
	"github.com/erigontech/mdbx-go/mdbx"
)

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

/*func (s *Schema) ExecuteQuery(query string) (interface{}, error) {
	lexer := parser.NewLexer(query)
	p := parser.NewParser(lexer)
	stmt := p.Parse().(*parser.SelectStatement)

	// Fetch data from the main table
	mainTable, err := s.GetTable(stmt.TableName)
	if err != nil {
		return nil, err
	}
	mainData, err := mainTable.fetchData(stmt.Conditions)
	if err != nil {
		return nil, err
	}

	// Apply joins
	for _, join := range stmt.Joins {
		joinTable, err := s.GetTable(join.Table)
		if err != nil {
			return nil, err
		}
		joinData, err := joinTable.fetchData(nil)
		if err != nil {
			return nil, err
		}
		mainData = applyJoin(mainData, joinData, join.OnField1, join.OnField2)
	}

	// Apply aggregates
	results := mainData
	for _, aggregate := range stmt.Aggregates {
		results = applyAggregate(results, aggregate)
	}

	return results, nil
}
*/
