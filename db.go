package mdbxsql

import (
	"context"
	"github.com/erigontech/mdbx-go/mdbx"
)

type Db struct {
	ctx context.Context
	env *mdbx.Env
	dbi mdbx.DBI
}

func NewDb(ctx context.Context, path string, maxDBs uint64) (*Db, error) {
	env, err := mdbx.NewEnv()
	if err != nil {
		return nil, err
	}

	err = env.SetGeometry(-1, -1, 1024*1024*1024*1024, -1, -1, 8192)
	if err != nil {
		return nil, err
	}

	err = env.SetOption(mdbx.OptMaxDB, maxDBs)
	if err != nil {
		return nil, err
	}

	err = env.Open(path, 0, 0664)
	if err != nil {
		return nil, err
	}

	var dbi mdbx.DBI
	err = env.Update(func(txn *mdbx.Txn) (err error) {
		dbi, err = txn.OpenRoot(mdbx.Create)
		return err
	})

	if err != nil {
		env.Close()
		return nil, err
	}

	return &Db{ctx: ctx, env: env, dbi: dbi}, nil
}

func (db *Db) GetEnv() *mdbx.Env {
	return db.env
}

func (db *Db) GetDBI() mdbx.DBI {
	return db.dbi
}

func (db *Db) Close() error {
	db.env.Close()
	return nil
}
