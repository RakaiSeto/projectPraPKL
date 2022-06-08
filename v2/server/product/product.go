package product

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/RakaiSeto/projectPraPKL/v2/db"
	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var dbconn *sql.DB
var varError error

func init() {
	dbconn = db.Db
}

func AllProduct() ([]*proto.Product, error) {
	rows, err := dbconn.Query("SELECT id, name, description, price FROM public.product")
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	defer rows.Close()

	products := make([]*proto.Product, 0)
	for rows.Next() {
		var product proto.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return nil, status.Error(codes.Code(2), err.Error())
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
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &product, nil
}

func AddProduct(product *proto.AdminProduct) (*proto.AddProductStatus, error) {
	userRow := dbconn.QueryRow("SELECT password, role FROM public.user where id=$1", product.GetAdminid())
	var i string
	var role string
	err := userRow.Scan(&i, &role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != product.GetAdminpass() {
		fmt.Println(i)
		fmt.Println(product.GetAdminpass())
		if product.GetAdminpass() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}
	if role != "admin" {
		return nil, status.Error(codes.Code(7), "please use admin account for this request")
	}
	
	var j int
	row := dbconn.QueryRow("SELECT id FROM public.product WHERE name=$1", product.GetName())
	_ = row.Scan(&j)
	if j != 0 {
		return nil, status.Error(codes.Code(6), "product already exists, please change its name")
	}
	
	_, err = dbconn.Exec("INSERT INTO public.product (name, description, price) VALUES ($1, $2, $3)", product.GetName(), product.GetDescription(), product.GetPrice())
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	row = dbconn.QueryRow("SELECT id FROM public.product WHERE name=$1", product.GetName())
	err = row.Scan(&product.Id)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	resp := proto.AddProductStatus{Response: "success", AdminProduct: product}

	return &resp, nil
}

func UpdateProduct(product *proto.AdminProduct) (*proto.ResponseStatus, error){
	userRow := dbconn.QueryRow("SELECT password, role FROM public.user where id=$1", product.GetAdminid())
	var i string
	var role string
	err := userRow.Scan(&i, &role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != product.GetAdminpass() {
		fmt.Println(i)
		fmt.Println(product.GetAdminpass())
		if product.GetAdminpass() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}
	if role != "admin" {
		return nil, status.Error(codes.Code(7), "please use admin account for this request")
	}
	
	QueryProduct := proto.Product{}
	row := dbconn.QueryRow("SELECT * from public.product where id = $1", product.Id)
	err = row.Scan(&QueryProduct.Id, &QueryProduct.Name, &QueryProduct.Description, &QueryProduct.Price)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "product not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}

	if product.Name != "" {QueryProduct.Name = product.Name}
	if product.Description != "" {QueryProduct.Description = product.Description}
	if product.Price != 0 {QueryProduct.Price = product.Price}

	_, err = dbconn.Exec("UPDATE public.product SET name=$2, description=$3, price=$4 WHERE id=$1", QueryProduct.Id, QueryProduct.Name, QueryProduct.Description, QueryProduct.Price)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}

func DeleteProduct(product *proto.AdminProduct) (*proto.ResponseStatus, error) {
	userRow := dbconn.QueryRow("SELECT password, role FROM public.user where id=$1", product.GetAdminid())
	var i string
	var role string
	err := userRow.Scan(&i, &role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != product.GetAdminpass() {
		fmt.Println(i)
		fmt.Println(product.GetAdminpass())
		if product.GetAdminpass() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}
	if role != "admin" {
		return nil, status.Error(codes.Code(7), "please use admin account for this request")
	}

	row := dbconn.QueryRow("SELECT name FROM public.product where id=$1", product.GetId())

	var name string
	
	err = row.Scan(&name)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	} else if name == "" {
		varError = fmt.Errorf("product not found")
		return nil, varError
	}

	_, err = dbconn.Exec("DELETE FROM public.product WHERE id=$1", product.GetId())
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}