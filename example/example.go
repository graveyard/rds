package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/Clever/rds"
)

var secretArn = os.Getenv("SECRET_ARN")
var resourceArn = os.Getenv("RESOURCE_ARN")

func main() {
	db, err := sql.Open("rds", fmt.Sprintf(`{
"SecretArn": "%s",
"ResourceArn": "%s"
}`, secretArn, resourceArn))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

}
