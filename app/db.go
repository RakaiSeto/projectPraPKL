package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Order struct{
	Id int
	Uname string
	Price int
	Profit int
}

type ProductOrder struct{
	Id int
	Orderid int
	Prodid int
	Qty int
	Beforediscount int
	Otherdisc int
	Poprice int
}

type Product struct {
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