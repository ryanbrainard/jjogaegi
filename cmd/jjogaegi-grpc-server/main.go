package main

import (
	"github.com/ryanbrainard/jjogaegi/grpc/proto"
	"github.com/ryanbrainard/jjogaegi/grpc/server"
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
	proto.RegisterJjogaegiServer(s, server.NewServer())
	s.Serve(listener)
}
