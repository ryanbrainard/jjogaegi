package server_test

import (
	"bytes"
	"context"
	"github.com/ryanbrainard/jjogaegi/grpc/proto"
	"github.com/ryanbrainard/jjogaegi/grpc/server"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"net"
	"testing"
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(1000)

	go func() {
		s := grpc.NewServer()
		proto.RegisterJjogaegiServer(s, server.NewServer())
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func newClient() (proto.JjogaegiClient, func()) {
	dialer := grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	})

	conn, err := grpc.DialContext(context.Background(), "", dialer, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return proto.NewJjogaegiClient(conn), func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}
}

func TestRun(t *testing.T) {
	c, teardown := newClient()
	defer teardown()

	response, err := c.Run(context.Background(), &proto.RunRequest{
		Config: &proto.RunConfig{
			Parser:    "list",
			Formatter: "csv",
		},
		Input: []byte("안녕 hello\n고양이 cat"),
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, string(response.Output), ",,안녕,,,hello,,,,,,,,,")
	assert.Contains(t, string(response.Output), ",,고양이,,,cat,,,,,,,,,")
}

func TestRunWithOptions(t *testing.T) {
	c, teardown := newClient()
	defer teardown()

	response, err := c.Run(context.Background(), &proto.RunRequest{
		Config: &proto.RunConfig{
			Parser:    "list",
			Formatter: "csv",
			Options: map[string]string{
				pkg.OPT_HEADER: "HEADER",
			},
		},
		Input: []byte("안녕 hello\n고양이 cat"),
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, string(response.Output), "HEADER")
	assert.Contains(t, string(response.Output), ",,안녕,,,hello,,,,,,,,,")
	assert.Contains(t, string(response.Output), ",,고양이,,,cat,,,,,,,,,")
}

func TestRunStream(t *testing.T) {
	c, teardown := newClient()
	defer teardown()

	stream, err := c.RunStream(context.TODO())
	if err != nil {
		panic(err)
	}
	waitc := make(chan struct{})
	out := &bytes.Buffer{}

	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(waitc)
					return
				}
				t.Fatal(err)
			}
			if _, err := out.Write(response.Output); err != nil {
				t.Fatal(err)
			}
		}
	}()

	req0 := &proto.RunRequest{
		Config: &proto.RunConfig{
			Parser:    "list",
			Formatter: "csv",
		},
		Input: []byte("안녕 hello\n"),
	}
	if err := stream.Send(req0); err != nil {
		t.Fatal(err)
	}

	req1 := &proto.RunRequest{
		Input: []byte("고양이 cat\n"),
	}
	if err := stream.Send(req1); err != nil {
		t.Fatal(err)
	}

	if err := stream.CloseSend(); err != nil {
		t.Fatal(err)
	}

	<-waitc

	assert.Contains(t, out.String(), ",,안녕,,,hello,,,,,,,,,")
	assert.Contains(t, out.String(), ",,고양이,,,cat,,,,,,,,,")
}
