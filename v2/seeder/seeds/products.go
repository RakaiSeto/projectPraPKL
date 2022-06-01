package seeds

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"
)

// ProductSeed seeds product data
func (s Seed) Product() {
	// check ada isinya ga
	row := s.db.QueryRow("SELECT id FROM product LIMIT 1")
	
	var i int
	err := row.Scan(&i)
	if err != nil {
		fmt.Println(err)
	}
	s.db.Exec("ALTER SEQUENCE product_id_seq RESTART")
	if i == 1{
		// klo ada isi baru di TRUNCATE
		s.db.Exec("TRUNCATE product")
	}

	for i := 0; i < 5; i++ {
		//prepare the statement
		stmt, err := s.db.Prepare(`INSERT INTO "product" (name, description, price) VALUES ($1, $2, $3)`)
		if err != nil {
			panic(err)
		}
		randInt, _ := faker.RandomInt(10, 150, 1)
		// execute query
		_, err = stmt.Exec(faker.Word(), faker.Sentence(), randInt[0])
		if err != nil {
			panic(err)
		}
	}
}