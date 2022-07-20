package data

import (
	"database/sql"
	"runtime"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	prom "github.com/prometheus/client_golang/prometheus"

	"github.rakops.com/rm/signal-api/stdlib/stats"
	"github.rakops.com/rm/signal-api/stdlib/stats/prometheus"
)

type (
	sqlxWrapperStats struct {
		db     *sqlx.DB
		s      stats.Client
		dbName string
	}

	txWrapperStats struct {
		tx     *sqlx.Tx
		s      stats.Client
		dbName string
	}
)

func labels(driver, operation, dbName string) stats.Labels {
	caller := "unknown"

	if pc, _, _, ok := runtime.Caller(3); ok {
		details := runtime.FuncForPC(pc)
		if details != nil {
			splitStr := strings.SplitAfter(details.Name(), ".")
			caller = splitStr[len(splitStr)-1]
		}
	}

	return stats.Labels{
		"driver", driver,
		"operation", operation,
		"db", dbName,
		"caller", caller,
	}
}

// NewSqlxWrapperStats returns a new instance that records stats to the provided client.
func NewSqlxWrapperStats(db *sqlx.DB, s stats.Client, dbName string) SqlxWrapper {
	return &sqlxWrapperStats{db: db, s: s, dbName: dbName}
}

// PrometheusCollectors is a prepopulated list of prometheus collectors.
// This must be used when using the stats collectors in this package
// with prometheus.
func PrometheusCollectors(app, team, env string) prometheus.Collectors {
	return prometheus.Collectors{
		"sql_operation": prom.NewSummaryVec(
			prom.SummaryOpts{
				Name: "sql_operation",
				Help: "The duration of the sql operation",
				ConstLabels: prom.Labels{
					"app":  app,
					"team": team,
					"env":  env,
				},
			},
			[]string{"driver", "operation", "db", "caller"},
		),
	}
}

func (s *sqlxWrapperStats) Get(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := s.db.Get(dest, query, args...)
	end := time.Now()

	s.s.Timing("sql_operation", s.labels("select"), end.Sub(start))

	return err
}

func (s *sqlxWrapperStats) Select(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := s.db.Select(dest, query, args...)
	end := time.Now()

	s.s.Timing("sql_operation", s.labels("select"), end.Sub(start))

	return err
}

func (s *sqlxWrapperStats) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := s.db.Query(query, args...)
	end := time.Now()

	s.s.Timing("sql_operation", s.labels("select"), end.Sub(start))

	return rows, err
}

func (s *sqlxWrapperStats) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := s.db.Exec(query, args...)
	end := time.Now()

	s.s.Timing("sql_operation", s.labels("exec"), end.Sub(start))

	return result, err
}

func (s *sqlxWrapperStats) NamedExec(query string, arg interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := s.db.NamedExec(query, arg)
	end := time.Now()

	s.s.Timing("sql_operation", s.labels("exec"), end.Sub(start))

	return result, err
}

func (s *sqlxWrapperStats) MustExec(query string, args ...interface{}) sql.Result {
	start := time.Now()
	result := s.db.MustExec(query, args...)
	end := time.Now()

	s.s.Timing("sql_operation", s.labels("exec"), end.Sub(start))

	return result
}

func (s *sqlxWrapperStats) Beginx() (TxWrapper, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &txWrapperStats{
		tx:     tx,
		s:      s.s,
		dbName: s.dbName,
	}, nil
}

func (s *sqlxWrapperStats) MustBegin() TxWrapper {
	return &txWrapperStats{
		tx:     s.db.MustBegin(),
		s:      s.s,
		dbName: s.dbName,
	}
}

func (s *sqlxWrapperStats) Rebind(query string) string {
	return s.db.Rebind(query)
}

func (s *sqlxWrapperStats) Stats() sql.DBStats {
	return s.db.Stats()
}

func (s *sqlxWrapperStats) DB() *sqlx.DB {
	return s.db
}

func (s *sqlxWrapperStats) labels(operation string) stats.Labels {
	return labels(s.db.DriverName(), operation, s.dbName)
}

func (t *txWrapperStats) Commit() error {
	start := time.Now()
	err := t.tx.Commit()
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("commit"), end.Sub(start))

	return err
}

func (t *txWrapperStats) Rollback() error {
	start := time.Now()
	err := t.tx.Rollback()
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("rollback"), end.Sub(start))

	return err
}

func (t *txWrapperStats) Get(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := t.tx.Get(dest, query, args...)
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("select"), end.Sub(start))

	return err
}

func (t *txWrapperStats) Select(dest interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := t.tx.Select(dest, query, args...)
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("select"), end.Sub(start))

	return err
}

func (t *txWrapperStats) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := t.tx.Query(query, args...)
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("select"), end.Sub(start))

	return rows, err
}

func (t *txWrapperStats) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := t.tx.Exec(query, args...)
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("exec"), end.Sub(start))

	return result, err
}

func (t *txWrapperStats) NamedExec(query string, arg interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := t.tx.NamedExec(query, arg)
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("exec"), end.Sub(start))

	return result, err
}

func (t *txWrapperStats) MustExec(query string, args ...interface{}) sql.Result {
	start := time.Now()
	result := t.tx.MustExec(query, args...)
	end := time.Now()

	t.s.Timing("sql_operation", t.labels("exec"), end.Sub(start))

	return result
}

func (t *txWrapperStats) Rebind(query string) string {
	return t.tx.Rebind(query)
}

func (t *txWrapperStats) labels(operation string) stats.Labels {
	return labels(t.tx.DriverName(), operation, t.dbName)
}
