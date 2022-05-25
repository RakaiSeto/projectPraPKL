package seeds

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"
)

func (s Seed) Testable() {

	for i := 0; i < 100; i++ {
		//prepare the statement
		stmt, err := s.db.Prepare("INSERT INTO testable(uname, email, password, role) VALUES ($1, $2, $3, $4)")
		if err != nil {
			panic(err)
		}
		// execute query
		_, err = stmt.Exec(fmt.Sprintln(faker.FirstName()+"123"), faker.Email(), "password", "customer")
		if err != nil {
			panic(err)
		}
	}
}
