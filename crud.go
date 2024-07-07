package mdbxsql

import (
	"encoding/json"
	"github.com/erigontech/mdbx-go/mdbx"
	"reflect"
)

func Insert(table *Table, record Model) error {
	primaryKey := reflect.ValueOf(record).Elem().FieldByName(table.Primary).Interface()
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	value, err := record.Marshal()
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Put(table.db.dbi, key, value, 0)
	})
}

func Update(table *Table, record Model) error {
	primaryKey := reflect.ValueOf(record).Elem().FieldByName(table.Primary).Interface()
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	value, err := record.Marshal()
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Put(table.db.dbi, key, value, 0)
	})
}

func Delete(table *Table, primaryKey interface{}) error {
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Del(table.db.dbi, key, nil)
	})
}

func Get(table *Table, primaryKey interface{}, record Model) error {
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	var value []byte
	err = table.db.env.View(func(txn *mdbx.Txn) error {
		value, err = txn.Get(table.db.dbi, key)
		return err
	})

	if err != nil {
		return err
	}

	return record.Unmarshal(value)
}
