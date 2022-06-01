package seeds

import (
	"fmt"

	faker "github.com/bxcodec/faker/v3"
)

func (s Seed) Order() {
	// check ada isinya ga
	row := s.db.QueryRow(`SELECT id FROM "order" LIMIT 1`)

	var i int
	err := row.Scan(&i)
	if err != nil {
		fmt.Println(err)
	}
	s.db.Exec("ALTER SEQUENCE order_id_seq RESTART")
	if i == 1 {
		// klo ada isi baru di TRUNCATE
		s.db.Exec("TRUNCATE public.order")
	}

	// how many id are there in user
	var userCount int
	s.db.QueryRow("SELECT count(*) FROM public.user").Scan(&userCount)
	// how many id are there in product
	var productCount int
	s.db.QueryRow("SELECT count(*) FROM product").Scan(&productCount)

	for i := 0; i < 5; i++ {
		randUser, _ := faker.RandomInt(1, userCount, 1)
		randProduct, _ := faker.RandomInt(1, productCount)
		// get product price
		var productPrice int
		s.db.QueryRow("SELECT price FROM public.product where id = $1", randProduct[0]).Scan(&productPrice)

		// random quantity
		randQty, _ := faker.RandomInt(1, 12, 1)

		// calculate price
		price := (productPrice * randQty[0])

		//prepare the statement
		stmt, err := s.db.Prepare("INSERT INTO public.order(userid, productid, quantity, totalprice) VALUES ($1, $2, $3, $4)")
		if err != nil {
			panic(err)
		}
		// execute query
		_, err = stmt.Exec(randUser[0], randProduct[0], randQty[0], price)
		if err != nil {
			panic(err)
		}
	}
}