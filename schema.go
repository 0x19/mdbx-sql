package mdbxsql

type Table struct {
	Name    string
	Primary string
	db      *Db
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

func (s *Schema) CreateTable(name string, primary string) *Table {
	table := &Table{
		Name:    name,
		Primary: primary,
		db:      s.db,
	}
	s.tables[name] = table
	return table
}
