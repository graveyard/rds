package rds

import (
	"context"
	"database/sql/driver"

	errors "golang.org/x/xerrors"
)

// stmt implements the driver.Stmt interfaces.
type stmt struct {
	ac    *conn
	query string
}

var _ driver.Stmt = &stmt{}
var _ driver.StmtExecContext = &stmt{}
var _ driver.StmtQueryContext = &stmt{}

// Close closes the statement.
func (as *stmt) Close() error {
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
func (as *stmt) NumInput() int {
	return 0 // TODO
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// Deprecated: Drivers should implement StmtExecContext instead (or additionally).
func (as *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.Errorf("deprecated: use ExecContext")
}

// ExecContext executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// ExecContext must honor the context timeout and return when it is canceled.
func (as *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return as.ac.ExecContext(ctx, as.query, args)
}

// Query executes a query that may return rows, such as a
// SELECT.
//
// Deprecated: Drivers should implement StmtQueryContext instead (or additionally).
func (as *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.Errorf("deprecated: use QueryContext")
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (as *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return as.ac.QueryContext(ctx, as.query, args)
}
