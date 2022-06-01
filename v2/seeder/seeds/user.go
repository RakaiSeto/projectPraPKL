package seeds

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"
)

func (s Seed) UserTrunc() {
	// check ada isinya ga
	row := s.db.QueryRow("SELECT id FROM public.user LIMIT 1")
	
	var i int
	err := row.Scan(&i)
	if err != nil {
		fmt.Println(err)
	}
	s.db.Exec("ALTER SEQUENCE user_id_seq RESTART")
	if i == 1{
		// klo ada isi baru di TRUNCATE
		s.db.Exec("TRUNCATE public.user")
		fmt.Println("truncated")
	}

	for i := 0; i < 10; i++ {
		//prepare the statement
		stmt, err := s.db.Prepare(`INSERT INTO "user" (uname, email, password, role) VALUES ($1, $2, $3, $4)`)
		if err != nil {
			panic(err)
		}
		// execute query
		_, err = stmt.Exec(fmt.Sprintf("%v%v", faker.FirstName(), "123"), faker.Email(), "password", "customer")
		if err != nil {
			panic(err)
		}
	}

}
func (s Seed) UserNoTrunc() {
	for i := 0; i < 5; i++ {
		//prepare the statement
		stmt, err := s.db.Prepare(`INSERT INTO "user" (uname, email, password, role) VALUES ($1, $2, $3, $4)`)
		if err != nil {
			panic(err)
		}
		// execute query
		_, err = stmt.Exec(fmt.Sprintf("%v%v", faker.FirstName(), "123"), faker.Email(), "password", "customer")
		if err != nil {
			panic(err)
		}
	}
}
