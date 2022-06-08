package order

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/RakaiSeto/projectPraPKL/v2/db"
	proto "github.com/RakaiSeto/projectPraPKL/v2/server/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var dbconn *sql.DB

func init() {
	dbconn = db.Db
}



func AllOrder(userInput *proto.User) ([]*proto.Order, error) {
	userRow := dbconn.QueryRow("SELECT password FROM public.user where id=$1", userInput.GetId())
	var i string
	err := userRow.Scan(&i)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != userInput.GetPassword() {
		fmt.Println(i)
		fmt.Println(userInput.GetPassword())
		if userInput.GetPassword() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}

	// get the order
	ordersRows, err := dbconn.Query("SELECT * FROM public.order where userid = $1 ORDER BY id", userInput.GetId())
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}
	defer ordersRows.Close()

	orders := make([]*proto.Order, 0)
	for ordersRows.Next() {
		var order proto.Order
		err := ordersRows.Scan(&order.Id, &order.Userid, &order.Productid, &order.Quantity, &order.Totalprice)
		if err != nil {
			return nil, status.Error(codes.Code(2), err.Error())
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func OneOrder(input *proto.Order) (*proto.Order, error) {
	userRow := dbconn.QueryRow("SELECT password FROM public.user where id=$1", input.GetUserid())
	var i string
	err := userRow.Scan(&i)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != input.GetUserpassword() {
		fmt.Println(i)
		fmt.Println(input.GetUserpassword())
		if input.GetUserpassword() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}

	// get the order
	orderRow := dbconn.QueryRow("SELECT * FROM public.order WHERE userid = $1 AND id = $2", input.GetUserid(), input.GetId())

	err = orderRow.Scan(&input.Id, &input.Userid, &input.Productid, &input.Quantity, &input.Totalprice)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}	
	
	empty := ""
	input.Userpassword = &empty 
	return input, nil
}

func AddOrder(orderInput *proto.Order) (*proto.AddOrderStatus, error) {
	userRow := dbconn.QueryRow("SELECT password FROM public.user where id=$1", orderInput.GetUserid())
	var i string
	err := userRow.Scan(&i)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != orderInput.GetUserpassword() {
		fmt.Println(i)
		fmt.Println(orderInput.GetUserpassword())
		if orderInput.GetUserpassword() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}
	
	// get the product
	row := dbconn.QueryRow("SELECT * FROM public.product where id=$1", orderInput.GetProductid())
	product := proto.Product{}
	err = row.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			
		} else {
			return nil, status.Error(codes.Code(2), err.Error())
		}
	}
	
	// input order
	totalprice := orderInput.GetQuantity() * product.Price
	_, err = dbconn.Exec("INSERT INTO public.order (userid, productid, quantity, totalprice) VALUES ($1, $2, $3, $4)", orderInput.GetUserid(), orderInput.GetProductid(), orderInput.GetQuantity(), totalprice)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	// get inputted order
	row = dbconn.QueryRow("SELECT * FROM public.order WHERE userid=$1 ORDER BY id desc limit 1", orderInput.GetUserid())

	err = row.Scan(&orderInput.Id, &orderInput.Userid, &orderInput.Productid, &orderInput.Quantity, &orderInput.Totalprice)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	empty := ""
	orderInput.Userpassword = &empty 
	return &proto.AddOrderStatus{Response: "success", Order: orderInput}, nil
}

func UpdateOrder(orderInput *proto.Order) (*proto.ResponseStatus, error) {
	userRow := dbconn.QueryRow("SELECT password FROM public.user where id=$1", orderInput.GetUserid())
	var i string
	err := userRow.Scan(&i)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != orderInput.GetUserpassword() {
		fmt.Println(i)
		fmt.Println(orderInput.GetUserpassword())
		if orderInput.GetUserpassword() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}
	
	// get the product
	row := dbconn.QueryRow("SELECT * FROM public.product where id=$1", orderInput.GetProductid())
	product := proto.Product{}
	err = row.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "product not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}

	// update order
	totalprice := orderInput.GetQuantity() * product.Price
	_, err = dbconn.Exec("UPDATE public.order SET userid=$1, productid=$2, quantity=$3, totalprice=$4 WHERE id=$5", orderInput.GetUserid(), orderInput.GetProductid(), orderInput.GetQuantity(), totalprice, orderInput.GetId())
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &proto.ResponseStatus{Response: "success"}, nil
}

func DeleteOrder(orderInput *proto.Order) (*proto.ResponseStatus, error) {
	userRow := dbconn.QueryRow("SELECT password FROM public.user where id=$1", orderInput.GetUserid())
	var i string
	err := userRow.Scan(&i)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	if i != orderInput.GetUserpassword() {
		fmt.Println(i)
		fmt.Println(orderInput.GetUserpassword())
		return nil, errors.New("wrong password for user")
	}

	// delete order
	_, err = dbconn.Exec("DELETE FROM public.order WHERE id=$1", orderInput.GetId())
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &proto.ResponseStatus{Response: "success"}, nil
}