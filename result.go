package rds

import (
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	errors "golang.org/x/xerrors"
)

// result implements the driver.Result interface.
type result struct {
	out *rdsdataservice.ExecuteStatementOutput
}

var _ driver.Result = &result{}

// LastInsertId returns the database's auto-generated ID
// after, for example, an INSERT into a table with primary
// key.
func (ar *result) LastInsertId() (int64, error) {
	if l := len(ar.out.GeneratedFields); l == 0 {
		return 0, errors.Errorf("no generated fields in result")
	} else if l != 1 {
		return 0, errors.Errorf("%d generated fields in result: %v", l, ar.out.GeneratedFields)
	}
	f := ar.out.GeneratedFields[0]
	if f.LongValue != nil {
		return aws.Int64Value(f.LongValue), nil
	}
	return 0, errors.Errorf("unhandled generated field type: %v", f)
}

// RowsAffected returns the number of rows affected by the
// query.
func (ar *result) RowsAffected() (int64, error) {
	return aws.Int64Value(ar.out.NumberOfRecordsUpdated), nil
}
