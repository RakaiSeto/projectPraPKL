package product

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RakaiSeto/projectPraPKL/v2/db"
	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
)

var dbconn *sql.DB
var varError error

func init() {
	dbconn = db.Db
}

func AllProduct() ([]*proto.Product, error) {
	rows, err := dbconn.Query("SELECT id, name, description, price FROM public.product")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	products := make([]*proto.Product, 0)
	for rows.Next() {
		var product proto.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func OneProduct(id int) (*proto.Product, error) {
	row := dbconn.QueryRow("SELECT id, name, description, price FROM public.product where id=$1", id)

	product := proto.Product{}
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func AddProduct(product *proto.Product) (*proto.AddProductStatus, error) {
	row := dbconn.QueryRow("SELECT id FROM public.product WHERE name=$1", product.GetName())
	var i int
	err := row.Scan(&i)
	if i != 0 {
		return nil, err
	}
	
	_, err = dbconn.Exec("INSERT INTO public.product (name, description, price) VALUES ($1, $2, $3)", product.GetName(), product.GetDescription(), product.GetPrice())
	if err != nil {
		errString := err.Error()
		resp := proto.AddProductStatus{Response: "failed", Error: &errString}
		return &resp, sql.ErrConnDone
	}

	row = dbconn.QueryRow("SELECT id FROM public.product WHERE name=$1", product.GetName())
	err = row.Scan(&product.Id)
	if err != nil {
		errString := err.Error()
		resp := proto.AddProductStatus{Response: "failed", Error: &errString}
		return &resp, sql.ErrConnDone
	}

	resp := proto.AddProductStatus{Response: "success", Product: product}

	return &resp, nil
}

func UpdateProduct(product *proto.Product) (*proto.ResponseStatus, error){
	QueryProduct := proto.Product{}
	row := dbconn.QueryRow("SELECT * from public.product where id = $1", product.Id)
	err := row.Scan(&QueryProduct.Id, &QueryProduct.Name, &QueryProduct.Description, &QueryProduct.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if product.Name != "" {QueryProduct.Name = product.Name}
	if product.Description != "" {QueryProduct.Description = product.Description}
	if product.Price != 0 {QueryProduct.Price = product.Price}

	_, err = dbconn.Exec("UPDATE public.product SET name=$2, description=$3, price=$4 WHERE id=$1", QueryProduct.Id, QueryProduct.Name, QueryProduct.Description, QueryProduct.Price)
	if err != nil {
		varError = err
		return nil, varError
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}

func DeleteProduct(id *proto.Id) (*proto.ResponseStatus, error) {
	row := dbconn.QueryRow("SELECT name FROM public.product where id=$1", id.Id)

	var name string
	
	err := row.Scan(&name)
	if err != nil {
		return nil, err
	} else if name == "" {
		varError = fmt.Errorf("user not found")
		return nil, varError
	}

	_, err = dbconn.Exec("DELETE FROM public.product WHERE id=$1", id.Id)
	if err != nil {
		return nil, err
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}