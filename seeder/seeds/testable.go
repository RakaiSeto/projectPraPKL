package seeds

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"
)

func (s Seed) Testable() {
	// check ada isinya ga
	row := s.db.QueryRow("SELECT id FROM testable LIMIT 1")
	
	var i int
	err := row.Scan(&i)
	if err != nil {
		panic(err)
	}
	s.db.Exec("ALTER SEQUENCE testable_id_seq RESTART")
	if i == 1{
		// klo ada isi baru di TRUNCATE
		s.db.Exec("TRUNCATE testable")
	}

	for i := 0; i < 15; i++ {
		//prepare the statement
		stmt, err := s.db.Prepare("INSERT INTO testable(uname, email, password, role) VALUES ($1, $2, $3, $4)")
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
