package main

import (
	"context"
	"github.com/ryanbrainard/jjogaegi/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := proto.NewJjogaegiClient(conn)

	response, err := c.Run(context.TODO(), &proto.RunRequest{
		Input:     "안녕 hello\n고양이 cat",
		Parser:    "list",
		Formatter: "json",
	})
	if err != nil {
		panic(err)
	}
	log.Printf("output \n%+v", response.Output)
}

type client struct {
}

func newClient() proto.JjogaegiClient {
	return &client{}
}

func (c *client) Run(ctx context.Context, in *proto.RunRequest, opts ...grpc.CallOption) (*proto.RunResponse, error) {
	return &proto.RunResponse{}, nil
}
