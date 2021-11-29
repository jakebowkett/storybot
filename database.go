package storybot

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type NullInt64 int64

func (n *NullInt64) Scan(src interface{}) error {
	switch m := src.(type) {
	case int64:
		*n = NullInt64(m)
	case nil:
		*n = 0
	default:
		return errors.New("invalid src for NullInt64")
	}
	return nil
}

type DbSliceStr = pq.StringArray

type DB struct {
	*sqlx.DB
}

func (db *DB) Beginx() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}
func (db *DB) Insert(tbl string, src interface{}, col []string) error {
	return insert(db.DB, tbl, src, col)
}
func (db *DB) Exists(fromWhere string, args ...interface{}) (bool, error) {
	return exists(db.DB, fromWhere, args...)
}

type Tx struct {
	*sqlx.Tx
}

func (tx *Tx) Insert(tbl string, src interface{}, col []string) error {
	return insert(tx.Tx, tbl, src, col)
}
func (tx *Tx) Exists(fromWhere string, args ...interface{}) (bool, error) {
	return exists(tx.Tx, fromWhere, args...)
}

type DBNamedExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}
type DBGetter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

func insert(x DBNamedExecer, tbl string, src interface{}, col []string) error {
	colNames := strings.Join(col, ", ")
	for i := range col {
		col[i] = ":" + col[i]
	}
	bindVars := strings.Join(col, ", ")
	_, err := x.NamedExec(fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		tbl,
		colNames,
		bindVars,
	), src)
	return err
}

func exists(x DBGetter, fromWhere string, args ...interface{}) (bool, error) {
	exists := false
	err := x.Get(
		&exists,
		fmt.Sprintf("SELECT EXISTS(SELECT 1 %s)", fromWhere),
		args...,
	)
	if err != nil {
		return false, err
	}
	return exists, nil
}
