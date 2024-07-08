package mdbxsql

import (
	"fmt"
	"github.com/erigontech/mdbx-go/mdbx"
	"time"
)

func BatchInsert(table *Table, records []Model) error {
	start := time.Now()
	return table.db.env.Update(func(txn *mdbx.Txn) error {
		for _, record := range records {
			value, err := record.Marshal()
			if err != nil {
				return err
			}
			if err := txn.Put(table.dbi, record.PrimaryKey(), value, 0); err != nil {
				return err
			}
		}
		fmt.Println("Time taken for batch operation:", time.Since(start))
		return nil
	})
}
