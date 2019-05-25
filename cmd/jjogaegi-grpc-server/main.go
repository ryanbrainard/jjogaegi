package main

import (
	"bytes"
	"context"
	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/proto"
	"github.com/ryanbrainard/jjogaegi/run"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
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

	output := &bytes.Buffer{}

	err := run.Run(
		strings.NewReader(req.Input),
		output,
		cmd.ParseOptParser(req.Parser),
		cmd.ParseOptFormatter(req.Formatter),
		map[string]string{},
	)

	return &proto.RunResponse{Output: output.String()}, err
}
