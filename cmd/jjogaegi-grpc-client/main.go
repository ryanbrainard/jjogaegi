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

	response, err := c.Run(context.TODO(), &proto.RunRequest{Ping: "hello"})
	if err != nil {
		panic(err)
	}
	log.Println("response", response.Pong)
}

type client struct {
}

func newClient() proto.JjogaegiClient {
	return &client{}
}

func (c *client) Run(ctx context.Context, in *proto.RunRequest, opts ...grpc.CallOption) (*proto.RunResponse, error) {
	return &proto.RunResponse{}, nil
}
