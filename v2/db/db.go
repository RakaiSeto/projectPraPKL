package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

var Db *sql.DB

func init() {
	var err error

    Db, err = sql.Open("postgres", "postgres://postgres:password@localhost/projectpkl?sslmode=disable")
    if err != nil {
        panic(err)
    }
	if err = Db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You've connected to the database")
}