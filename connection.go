package rds

import (
	"context"
	"database/sql/driver"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	errors "golang.org/x/xerrors"
)

//go:generate mockgen -package rds -source $PWD/vendor/github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface/interface.go -destination rdsdataservice_mocks_test.go RDSDataServiceAPI

// conn implements the database/sql connection interfaces.
type conn struct {
	rds         rdsdataserviceiface.RDSDataServiceAPI
	database    string
	resourceArn string
	secretArn   string
}

var _ driver.Conn = &conn{}
var _ driver.ConnPrepareContext = &conn{}
var _ driver.QueryerContext = &conn{}
var _ driver.ExecerContext = &conn{}

// Begin is TODO.
func (ac *conn) Begin() (driver.Tx, error) {
	panic("TODO: Begin")
}

// Prepare is deprecated.
func (ac *conn) Prepare(query string) (driver.Stmt, error) {
	panic("deprecated: use PrepareContext")
}

// PrepareContext prepares a query.
func (ac *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return &stmt{
		ac:    ac,
		query: query,
	}, nil
}

// Ping checks connectivity.
func (ac *conn) Ping(ctx context.Context) error {
	_, err := ac.rds.ExecuteStatementWithContext(ctx, &rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(ac.resourceArn),
		Database:    aws.String(ac.database),
		SecretArn:   aws.String(ac.secretArn),
		Sql:         aws.String("/* ping */ SELECT 1"),
	})
	return err
}

// QueryContext carries out a basic SQL query.
func (ac *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "?") {
		// TODO: support ordinal args
		return nil, errors.Errorf("ordinal parameters not supported, must used named parameters")
	}

	parameters := make([]*rdsdataservice.SqlParameter, len(args))
	for i := range args {
		name := args[i].Name
		if name == "" {
			return nil, errors.Errorf("only named parameters supported")
		}
		value := args[i].Value
		if value == nil {
			parameters[i] = &rdsdataservice.SqlParameter{
				Name:  &name,
				Value: &rdsdataservice.Field{IsNull: aws.Bool(true)},
			}
			continue
		}
		var f *rdsdataservice.Field
		switch t := value.(type) {
		case string:
			f = &rdsdataservice.Field{StringValue: aws.String(t)}
		case []byte:
			f = &rdsdataservice.Field{BlobValue: t}
		case bool:
			f = &rdsdataservice.Field{BooleanValue: &t}
		case float64:
			f = &rdsdataservice.Field{DoubleValue: &t}
		case int64:
			f = &rdsdataservice.Field{LongValue: &t}
		default:
			return nil, errors.Errorf("%s is unsupported type: %#v", name, value)
		}
		parameters[i] = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: f,
		}
	}

	out, err := ac.rds.ExecuteStatementWithContext(ctx, &rdsdataservice.ExecuteStatementInput{
		ContinueAfterTimeout:  aws.Bool(false),
		ResourceArn:           aws.String(ac.resourceArn),
		Database:              aws.String(ac.database),
		IncludeResultMetadata: aws.Bool(true),
		SecretArn:             aws.String(ac.secretArn),
		Sql:                   aws.String(query),
		Parameters:            parameters,
	})
	if err != nil {
		return nil, errors.Errorf("ExecuteStatement: %v", err)
	}
	return &rows{out: out}, nil
}

// Exec is deprecated, use ExecContext.
func (ac *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return nil, errors.Errorf("deprecated: use ExecContext")
}

// ExecContext performs a query that doesn't return results, e.g. inserts.
func (ac *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	rs, err := ac.QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return &result{out: rs.(*rows).out}, nil
}

// Close the connection is a no-op.
func (ac *conn) Close() (err error) {
	ac.rds = nil // set this to nil to trigger garbage collection
	return nil
}
