package sql

import (
  "errors"
)

type Result struct {
  RowCount int64 `json:"rowcount"`
  // Duration int   `json:"duration"`
}

// LastInsertId returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary
// key.
func (result *Result) LastInsertId() (int64, error) {
  return 0, errors.New("sql: crate driver does not support last insert id")
}

// RowsAffected returns the number of rows affected by the
// query.
func (result *Result) RowsAffected() (int64, error) {
  return result.RowCount, nil
}
