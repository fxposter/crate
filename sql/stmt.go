package sql

import (
  "bytes"
  "database/sql/driver"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "regexp"
)

type Stmt struct {
  conn  *Conn
  query string
}

// Close closes the statement.
//
// As of Go 1.1, a Stmt will not be closed if it's in use
// by any queries.
func (stmt *Stmt) Close() error {
  return nil
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
var placeholderRegexp = regexp.MustCompile("\\$\\d+")

func (stmt *Stmt) NumInput() int {
  return len(placeholderRegexp.FindAll([]byte(stmt.query), -1))
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
  attributes := map[string]interface{}{
    "stmt": stmt.query,
  }
  if len(args) > 0 {
    attributes["args"] = args
  }

  body, err := json.Marshal(attributes)
  log.Println(string(body))
  if err != nil {
    return nil, err
  }
  server := stmt.conn.server()
  resp, err := http.Post(server+"/_sql", "application/json", bytes.NewReader(body))
  // TODO handle HTTP errors
  if err != nil {
    return nil, err
  }

  buffer, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  result := Result{}
  err = json.Unmarshal(buffer, &result)
  if err != nil {
    return nil, err
  }
  return &result, nil
}

// Query executes a query that may return rows, such as a
// SELECT.
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
  attributes := map[string]interface{}{
    "stmt": stmt.query,
  }
  if len(args) > 0 {
    attributes["args"] = args
  }

  body, err := json.Marshal(attributes)
  if err != nil {
    return nil, err
  }
  server := stmt.conn.server()
  resp, err := http.Post(server+"/_sql", "application/json", bytes.NewReader(body))
  if err != nil {
    return nil, err
  }

  buffer, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  result := Rows{currentRowIndex: 0}
  err = json.Unmarshal(buffer, &result)
  if err != nil {
    return nil, err
  }
  return &result, nil
}
