package rds

import (
	"database/sql/driver"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	errors "golang.org/x/xerrors"
)

// rows implements the driver.Rows interface.
type rows struct {
	out            *rdsdataservice.ExecuteStatementOutput
	recordPosition int
}

var _ driver.Rows = &rows{}

// By default we return a single column, which embodies the entire row response from the query.
func (r *rows) Columns() []string {
	cols := make([]string, len(r.out.ColumnMetadata))
	for i, c := range r.out.ColumnMetadata {
		cols[i] = aws.StringValue(c.Name)
	}
	return cols
}

func (r *rows) Next(dest []driver.Value) error {
	if r.recordPosition == len(r.out.Records) {
		return io.EOF
	}
	row := r.out.Records[r.recordPosition]
	r.recordPosition++

	for i, field := range row {
		coerced, err := convertField(field)
		if err != nil {
			return errors.Errorf("convertValue(col=%d): %v", i, err)
		}
		dest[i] = coerced
	}

	return nil
}

func convertField(field *rdsdataservice.Field) (interface{}, error) {
	switch {
	case field.BlobValue != nil:
		return field.BlobValue, nil
	case field.BooleanValue != nil:
		return *field.BooleanValue, nil
	case field.DoubleValue != nil:
		return *field.DoubleValue, nil
	case field.IsNull != nil:
		return nil, nil
	case field.LongValue != nil:
		return *field.LongValue, nil
	case field.StringValue != nil:
		return *field.StringValue, nil
	default:
		return nil, errors.Errorf("no part of Field non-nil")
	}
}

func (r *rows) Close() error {
	return nil
}
