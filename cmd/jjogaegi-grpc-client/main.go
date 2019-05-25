package main

import (
	"context"
	jgrpc "github.com/ryanbrainard/jjogaegi/grpc"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := jgrpc.NewJjogaegiClient(conn)

	response, err := c.Run(context.TODO(), &jgrpc.RunRequest{Ping: "hello"})
	if err != nil {
		panic(err)
	}
	log.Println("response", response.Pong)
}

type client struct {
}

func newClient() jgrpc.JjogaegiClient {
	return &client{}
}

func (c *client) Run(ctx context.Context, in *jgrpc.RunRequest, opts ...grpc.CallOption) (*jgrpc.RunResponse, error) {
	return &jgrpc.RunResponse{}, nil
}
