package data

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type (
	// DataContext contains basic database operations shared by basic and transaction level sqlx requests
	DataContext interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		NamedExec(query string, arg interface{}) (sql.Result, error)
		MustExec(query string, args ...interface{}) sql.Result
		Get(dest interface{}, query string, args ...interface{}) error
		Select(dest interface{}, query string, args ...interface{}) error
		Query(query string, args ...interface{}) (*sql.Rows, error)
		Rebind(query string) string
	}
	// SqlxWrapper composes DataContext and adds other methods.
	SqlxWrapper interface {
		DataContext
		Beginx() (TxWrapper, error)
		MustBegin() TxWrapper
		Stats() sql.DBStats
		DB() *sqlx.DB
	}
	// TxWrapper composes DataContext and adds transaction functions.
	TxWrapper interface {
		Commit() error
		Rollback() error
		DataContext
	}
	sqlxWrapperImpl struct {
		db *sqlx.DB
	}
	txWrapperImpl struct {
		tx *sqlx.Tx
	}
)

// NewSqlxWrapper returns a new plain instance.
func NewSqlxWrapper(db *sqlx.DB) SqlxWrapper {
	return &sqlxWrapperImpl{db: db}
}

// Basic  Implementation Details

func (s *sqlxWrapperImpl) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *sqlxWrapperImpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return s.db.NamedExec(query, arg)
}

func (s *sqlxWrapperImpl) MustExec(query string, args ...interface{}) sql.Result {
	return s.db.MustExec(query, args...)
}

func (s *sqlxWrapperImpl) Get(dest interface{}, query string, args ...interface{}) error {
	return s.db.Get(dest, query, args...)
}

func (s *sqlxWrapperImpl) Select(dest interface{}, query string, args ...interface{}) error {
	return s.db.Select(dest, query, args...)
}

func (s *sqlxWrapperImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *sqlxWrapperImpl) Beginx() (TxWrapper, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &txWrapperImpl{tx: tx}, nil
}

func (s *sqlxWrapperImpl) MustBegin() TxWrapper {
	return &txWrapperImpl{tx: s.db.MustBegin()}
}

func (s *sqlxWrapperImpl) Rebind(query string) string {
	return s.db.Rebind(query)
}

func (s *sqlxWrapperImpl) Stats() sql.DBStats {
	return s.db.Stats()
}

func (s *sqlxWrapperImpl) DB() *sqlx.DB {
	return s.db
}

// Transaction Implementation Details

func (t *txWrapperImpl) Commit() error {
	return t.tx.Commit()
}

func (t *txWrapperImpl) Rollback() error {
	return t.tx.Rollback()
}

func (t *txWrapperImpl) Get(dest interface{}, query string, args ...interface{}) error {
	return t.tx.Get(dest, query, args...)
}

func (t *txWrapperImpl) Select(dest interface{}, query string, args ...interface{}) error {
	return t.tx.Select(dest, query, args...)
}

func (t *txWrapperImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

func (t *txWrapperImpl) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *txWrapperImpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return t.tx.NamedExec(query, arg)
}

func (t *txWrapperImpl) MustExec(query string, args ...interface{}) sql.Result {
	return t.tx.MustExec(query, args...)
}

func (t *txWrapperImpl) Rebind(query string) string {
	return t.tx.Rebind(query)
}
