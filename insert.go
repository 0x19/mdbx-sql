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

func BatchInsert(table *Table, records []Model) error {
	start := time.Now()
	defer func() {
		fmt.Println("Time taken for batch operation:", time.Since(start))
	}()

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		cursor, err := txn.OpenCursor(table.dbi)
		if err != nil {
			return err
		}
		defer cursor.Close()

		for _, record := range records {
			value, err := record.Marshal()
			if err != nil {
				return err
			}
			if err := cursor.Put(record.PrimaryKey(), value, 0); err != nil {
				return err
			}
		}

		return nil
	})
}
