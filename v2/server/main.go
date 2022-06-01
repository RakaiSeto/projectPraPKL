package main

import (
	"fmt"
	"net"

	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	fmt.Println("1")
	
	srv := grpc.NewServer()
	fmt.Println("2")
	proto.RegisterServiceServer(srv, &Server{})
	fmt.Println("3")
	reflection.Register(srv)
	fmt.Println("4")

	if e := srv.Serve(listener); e != nil {
			panic(e)
		} else {
			fmt.Println("Server listening")
		}
}
