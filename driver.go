// Package rds provides a driver for the RDS Data API.
package rds

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// rdsDriver implements the driver.Driver interface.
type rdsDriver struct {
}

var _ driver.Driver = rdsDriver{}

type config struct {
	ResourceArn string
	SecretArn   string
}

// Open a connection. Parse the URL as a JSON config object.
func (d rdsDriver) Open(url string) (driver.Conn, error) {
	var c config
	if err := json.Unmarshal([]byte(url), &c); err != nil {
		return nil, err
	}
	sess := session.New(&aws.Config{Region: aws.String("us-west-2")})
	rdsAPI := rdsdataservice.New(sess)
	return &conn{
		rds:         rdsAPI,
		resourceArn: c.ResourceArn,
		secretArn:   c.SecretArn,
	}, nil
}

func init() {
	sql.Register("rds", &rdsDriver{})
}
