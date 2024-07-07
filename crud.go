package mdbxsql

import (
	"encoding/json"
	"github.com/erigontech/mdbx-go/mdbx"
	"reflect"
)

func Insert[T any](table *Table, record T) error {
	primaryKey := reflect.ValueOf(record).FieldByName(table.Primary).Interface()
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	value, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Put(table.db.dbi, key, value, 0)
	})
}

func Update[T any](table *Table, record T) error {
	primaryKey := reflect.ValueOf(record).FieldByName(table.Primary).Interface()
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	value, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Put(table.db.dbi, key, value, 0)
	})
}

func Delete[T any](table *Table, primaryKey T) error {
	key, err := json.Marshal(primaryKey)
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Del(table.db.dbi, key, nil)
	})
}

func Get[T any](table *Table, primaryKey any) (T, error) {
	key, err := json.Marshal(primaryKey)
	if err != nil {
		var zero T
		return zero, err
	}

	var value []byte
	err = table.db.env.View(func(txn *mdbx.Txn) error {
		value, err = txn.Get(table.db.dbi, key)
		return err
	})

	if err != nil {
		var zero T
		return zero, err
	}

	var record T
	err = json.Unmarshal(value, &record)
	if err != nil {
		var zero T
		return zero, err
	}

	return record, nil
}
