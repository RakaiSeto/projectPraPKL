package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/RakaiSeto/projectPraPKL/seeder/seeds"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	// "github.com/joho/godotenv"
)

func main() {
	handleArgs()
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			db, err := sql.Open("postgres", "postgres://postgres:password@localhost/prepkl?sslmode=disable")
			if err != nil {
				panic(err)
			}
		
			if err = db.Ping(); err != nil {
				panic(err)
			}
			seeds.Execute(db, args[1:]...)
			os.Exit(0)
		}
	}
}
