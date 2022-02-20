package config

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/prepkl?sslmode=false")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You've connected to the database")
}