package main

import (
	"context"
	"net"

	user "github.com/RakaiSeto/projectPraPKL/v2/server/user"
	order "github.com/RakaiSeto/projectPraPKL/v2/server/order"
	product "github.com/RakaiSeto/projectPraPKL/v2/server/product"
	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{
	proto.ServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	s := Server{}

	srv := grpc.NewServer()
	proto.RegisterServiceServer(srv, &s)
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}

func (s *Server) Tes(ctx context.Context, empty *proto.EmptyStruct) (*proto.ResponseStatus, error) {
	return &proto.ResponseStatus{Response: "Hello"}, nil
}

func (s *Server) AllUser(ctx context.Context, empty *proto.EmptyStruct) (*proto.Users, error) {
	response, err := user.AllUser()
	if err != nil {
		return nil, err
	}
	var returned proto.Users
	returned.User = response
	return &returned, nil 
}

func (s *Server) OneUser(ctx context.Context, id *proto.Id) (*proto.User, error) {
	response, err := user.OneUser(int(id.GetId()))
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) AddUser(ctx context.Context, userInput *proto.User) (*proto.AddUserStatus, error) {
	response, err := user.AddUser(userInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) UpdateUser(ctx context.Context, userInput *proto.User) (*proto.ResponseStatus, error) {
	response, err := user.UpdateUser(userInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) DeleteUser(ctx context.Context, userInput *proto.User) (*proto.ResponseStatus, error) {
	response, err := user.DeleteUser(userInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) AllProduct(ctx context.Context, empty *proto.EmptyStruct) (*proto.Products, error) {
	response, err := product.AllProduct()
	if err != nil {
		return nil, err
	}
	var returned proto.Products
	returned.Product = response
	return &returned, nil 
}

func (s *Server) OneProduct(ctx context.Context, id *proto.Id) (*proto.Product, error) {
	response, err := product.OneProduct(int(id.GetId()))
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) AddProduct(ctx context.Context, productInput *proto.AdminProduct) (*proto.AddProductStatus, error) {
	response, err := product.AddProduct(productInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) UpdateProduct(ctx context.Context, productInput *proto.AdminProduct) (*proto.ResponseStatus, error) {
	response, err := product.UpdateProduct(productInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) DeleteProduct(ctx context.Context, productInput *proto.AdminProduct) (*proto.ResponseStatus, error) {
	response, err := product.DeleteProduct(productInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) AllOrder(ctx context.Context, userInput *proto.User) (*proto.Orders, error) {
	response, err := order.AllOrder(userInput)
	if err != nil {
		return nil, err
	}
	var returned proto.Orders
	returned.Order = response
	return &returned, nil 
}

func (s *Server) OneOrder(ctx context.Context, orderInput *proto.Order) (*proto.Order, error) {
	response, err := order.OneOrder(orderInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) AddOrder(ctx context.Context, orderInput *proto.Order) (*proto.AddOrderStatus, error) {
	response, err := order.AddOrder(orderInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) UpdateOrder(ctx context.Context, orderInput *proto.Order) (*proto.ResponseStatus, error) {
	response, err := order.UpdateOrder(orderInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}

func (s *Server) DeleteOrder(ctx context.Context, orderInput *proto.Order) (*proto.ResponseStatus, error) {
	response, err := order.DeleteOrder(orderInput)
	if err != nil {
		return nil, err
	}
	return response, nil 
}