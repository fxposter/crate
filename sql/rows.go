package sql

import (
  "database/sql/driver"
  "io"
)

type Rows struct {
  Cols     []string         `json:"cols"`
  Rows     [][]driver.Value `json:"rows"`
  RowCount int64            `json:"rowcount"`
  Duration int              `json:"duration"`

  currentRowIndex int64
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (rows *Rows) Columns() []string {
  return rows.Cols
}

// Close closes the rows iterator.
func (rows *Rows) Close() error {
  return nil
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// The dest slice may be populated only with
// a driver Value type, but excluding string.
// All string values must be converted to []byte.
//
// Next should return io.EOF when there are no more rows.
func (rows *Rows) Next(dest []driver.Value) error {
  if rows.currentRowIndex >= rows.RowCount {
    return io.EOF
  } else {
    currentRow := rows.Rows[rows.currentRowIndex]
    for i, _ := range rows.Cols {
      dest[i] = currentRow[i]
    }
    rows.currentRowIndex++
    return nil
  }
}
