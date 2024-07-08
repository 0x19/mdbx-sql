package mdbxsql

import (
	"fmt"
	"github.com/erigontech/mdbx-go/mdbx"
	"time"
)

func Insert(table *Table, record Model) error {
	start := time.Now()
	value, err := record.Marshal()
	if err != nil {
		return err
	}
	fmt.Println("Time taken:", time.Since(start))

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
