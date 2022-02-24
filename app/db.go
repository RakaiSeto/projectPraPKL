package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type order struct{
	id int
	uname string
	price int
}

type productOrder struct{
	id int
	orderid int
	prodod int
	qty int
	otherdisc int
	poprice int
}

type product struct {
	id int
	prodcode int
	name string
	catprice int
	memprice int
	discount int
}

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/prepkl?sslmode=disable")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You've connected to the database")
}