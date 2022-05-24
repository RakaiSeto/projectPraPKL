package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

type Appuser struct {
	Id 		 int	`json:"id"`
	Uname    string `json:"uname"`
	Password string `json:"password"`
	Role 	 string `json:"role"`
}

type Order struct{
	Id int
	Uname string
	Price int
	Profit int
	Created string
	Role string
}

type Productorder struct{
	Id int
	Orderid int
	Prodcode int
	Qty int
	Discount int
	Poprice int
	Otherexp int
	Created string
	Otherdiscount int
	Profit int
	Role string
}

type Product struct {
	Prodcode int
	Name string
	Catprice int
}

var db *sql.DB
var rdb *redis.Client
var ctx = context.TODO()

func init() {
	// var ctx = context.Background()

	var err error

	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/prepkl?sslmode=disable")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You've connected to the database")

	rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	fmt.Println(rdb.Ping(ctx).Result())
}