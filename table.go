package mdbxsql

import (
	"github.com/erigontech/mdbx-go/mdbx"
)

type Table struct {
	Name    string
	Primary string
	db      *Db
	dbi     mdbx.DBI
}
