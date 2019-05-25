package main

import (
	"context"
	"github.com/ryanbrainard/jjogaegi/proto"
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
	proto.RegisterJjogaegiServer(s, newServer())
	s.Serve(listener)
}

type server struct {
}

func newServer() proto.JjogaegiServer {
	return &server{}
}

func (s *server) Run(_ context.Context, req *proto.RunRequest) (*proto.RunResponse, error) {
	return &proto.RunResponse{Pong: req.Ping + " world"}, nil
}
