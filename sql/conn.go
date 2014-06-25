package sql

import (
  "database/sql/driver"
  "errors"
)

type Conn struct {
  servers       []string
  currentServer int
}

func (conn *Conn) server() string {
  server := conn.servers[conn.currentServer]
  conn.currentServer = (conn.currentServer + 1) % len(conn.servers)
  return server
}

// Prepare returns a prepared statement, bound to this connection.
func (conn *Conn) Prepare(query string) (driver.Stmt, error) {
  return &Stmt{conn, query}, nil
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package cratetains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (conn *Conn) Close() error {
  return nil
}

// Begin starts and returns a new transaction.
func (conn *Conn) Begin() (driver.Tx, error) {
  return nil, errors.New("sql: crate driver does not support transactions")
}
