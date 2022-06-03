package main

import (
	"log"

	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var Client proto.ServiceClient
func init() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	Client = proto.NewServiceClient(conn)

}

func main() {
	g := gin.Default()
	g.GET("/hello", Tes)
	g.GET("/user", AllUser)
	g.GET("/user/:id", OneUser)
	g.POST("/user", PostUser)
	g.PATCH("/user/:id", PatchUser)
	g.DELETE("/user/:id", DeleteUser)

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}