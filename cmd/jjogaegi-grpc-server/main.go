package main

import (
	"context"
	jgrpc "github.com/ryanbrainard/jjogaegi/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5000})
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("listening network=%v string=%v", listener.Addr().Network(), listener.Addr().String())

	s := grpc.NewServer()
	jgrpc.RegisterJjogaegiServer(s, newServer())
	s.Serve(listener)
}

type server struct {
}

func newServer() jgrpc.JjogaegiServer {
	return &server{}
}

func (s *server) Run(_ context.Context, req *jgrpc.RunRequest) (*jgrpc.RunResponse, error) {
	return &jgrpc.RunResponse{Pong: req.Ping + " world"}, nil
}
