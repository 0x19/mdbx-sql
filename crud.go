package mdbxsql

import (
	"github.com/erigontech/mdbx-go/mdbx"
)

func Insert(table *Table, record Model) error {
	value, err := record.Marshal()
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Put(table.db.dbi, record.PrimaryKey(), value, 0)
	})
}

func Update(table *Table, record Model) error {
	primaryKey := record.PrimaryKey()

	value, err := record.Marshal()
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Put(table.db.dbi, primaryKey, value, 0)
	})
}

func Delete(table *Table, primaryKey []byte) error {
	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Del(table.db.dbi, primaryKey, nil)
	})
}

func Get(table *Table, primaryKey []byte, record Model) error {
	var value []byte
	var err error
	err = table.db.env.View(func(txn *mdbx.Txn) error {
		value, err = txn.Get(table.db.dbi, primaryKey)
		return err
	})

	if err != nil {
		return err
	}

	return record.Unmarshal(value)
}
