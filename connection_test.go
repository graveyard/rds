package rds

import (
	"context"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/golang/mock/gomock"
)

const mockDatabase = "mock-database"
const mockResourceArn = "arn:rds:mock"
const mockSecretArn = "arn:secret:mock"

func Test_conn_QueryContext(t *testing.T) {
	type fields struct {
		rds         func(*MockRDSDataServiceAPI) *MockRDSDataServiceAPI
		database    string
		resourceArn string
		secretArn   string
	}
	type args struct {
		ctx   context.Context
		query string
		args  []driver.NamedValue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    driver.Rows
		wantErr bool
	}{
		{
			name: "ordinal parameters are TODO",
			fields: fields{
				rds: func(m *MockRDSDataServiceAPI) *MockRDSDataServiceAPI {
					return m
				},
				database:    mockDatabase,
				resourceArn: mockResourceArn,
				secretArn:   mockSecretArn,
			},
			args: args{
				ctx:   context.Background(),
				query: "SELECT * FROM foo where x=?",
				args:  nil,
			},
			wantErr: true,
		},
		{
			name: "no parameters",
			fields: fields{
				rds: func(m *MockRDSDataServiceAPI) *MockRDSDataServiceAPI {
					m.EXPECT().ExecuteStatementWithContext(gomock.Any(), &rdsdataservice.ExecuteStatementInput{
						ContinueAfterTimeout:  aws.Bool(false),
						Database:              aws.String(mockDatabase),
						IncludeResultMetadata: aws.Bool(true),
						Parameters:            []*rdsdataservice.SqlParameter{},
						ResourceArn:           aws.String(mockResourceArn),
						Schema:                nil, // TODO might need to support parsing this out of the query?
						SecretArn:             aws.String(mockSecretArn),
						Sql:                   aws.String("SELECT * FROM foo"),
						TransactionId:         nil,
					})
					return m
				},
				database:    mockDatabase,
				resourceArn: mockResourceArn,
				secretArn:   mockSecretArn,
			},
			args: args{
				ctx:   context.Background(),
				query: "SELECT * FROM foo",
				args:  nil,
			},
			want:    &rows{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			rds := NewMockRDSDataServiceAPI(ctrl)
			ac := &conn{
				rds:         tt.fields.rds(rds),
				database:    tt.fields.database,
				resourceArn: tt.fields.resourceArn,
				secretArn:   tt.fields.secretArn,
			}
			got, err := ac.QueryContext(tt.args.ctx, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("conn.QueryContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("conn.QueryContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
