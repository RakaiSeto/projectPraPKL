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
	Created string
}

type Productorder struct{
	Id int
	Orderid int
	Procode int
	Qty int
	Discount int
	Poprice int
	Otherexp int
	Created string
	Otherdiscount int
	Profit int
}

type Product struct {
	Prodcode int
	Name string
	Catprice int
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