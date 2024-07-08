package mdbxsql

import "github.com/erigontech/mdbx-go/mdbx"

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
