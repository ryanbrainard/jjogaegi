package main

import (
	"github.com/ryanbrainard/jjogaegi/grpc/go/jjogaegigprc"
	"google.golang.org/grpc"
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

	addr := &net.TCPAddr{Port: portInt}
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := jjogaegigprc.NewRunServiceClient(conn)

	//runner := &jjogaegigprc.SimpleRunner{c}
	runner := &jjogaegigprc.StreamingRunner{c}

	err = runner.Run(os.Stdin, os.Stdout, &jjogaegigprc.RunConfig{
		Parser:    "list",
		Formatter: "csv",
	})
	if err != nil {
		panic(err)
	}
}
