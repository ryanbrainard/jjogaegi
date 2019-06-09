package main

import (
	"go.ryanbrainard.com/jjogaegi/grpc/go/jjogaegigprc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: portInt})
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("listening network=%v string=%v", listener.Addr().Network(), listener.Addr().String())

	s := grpc.NewServer()
	jjogaegigprc.RegisterRunServiceServer(s, jjogaegigprc.NewServer())
	s.Serve(listener)
}
