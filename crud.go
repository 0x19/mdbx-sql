package mdbxsql

import (
	"errors"
	"github.com/erigontech/mdbx-go/mdbx"
	"strconv"
)

func getKeyValue(primaryKey interface{}) ([]byte, error) {
	switch v := primaryKey.(type) {
	case int:
		return []byte(strconv.Itoa(v)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case string:
		return []byte(v), nil
	default:
		return nil, errors.New("unsupported key type")
	}
}

func Insert(table *Table, record Model) error {
	primaryKey := record.PrimaryKey()
	key, err := getKeyValue(primaryKey)
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
	primaryKey := record.PrimaryKey()
	key, err := getKeyValue(primaryKey)
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
	key, err := getKeyValue(primaryKey)
	if err != nil {
		return err
	}

	return table.db.env.Update(func(txn *mdbx.Txn) error {
		return txn.Del(table.db.dbi, key, nil)
	})
}

func Get(table *Table, primaryKey interface{}, record Model) error {
	key, err := getKeyValue(primaryKey)
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
