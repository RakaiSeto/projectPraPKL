package main

import (
	// "github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	_ = proto.NewServiceClient(conn)
}