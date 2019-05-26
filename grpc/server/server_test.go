package server_test

import (
	"context"
	"fmt"
	"github.com/ryanbrainard/jjogaegi/grpc/proto"
	"github.com/ryanbrainard/jjogaegi/grpc/server"
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

	response, err := c.Run(context.TODO(), &proto.RunRequest{
		Options: &proto.RunOptions{
			Parser:    "list",
			Formatter: "json",
		},
		Input: []byte("안녕 hello\n고양이 cat"),
	})
	if err != nil {
		panic(err)
	}
	t.Logf("output \n%+v", string(response.Output))
}

func TestRunStream(t *testing.T) {
	c, teardown := newClient()
	defer teardown()

	stream, err := c.RunStream(context.TODO())
	if err != nil {
		panic(err)
	}
	waitc := make(chan struct{})

	go func() {
		for {
			t.Log("waiting for response")
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					t.Log("waitc.close")
					close(waitc)
					return
				}
				panic(fmt.Sprintf("-> %+v", err))
			}
			t.Logf("output (%d) \n%+v", len(response.Output), string(response.Output))
		}
	}()

	t.Log("stream.Send[0]")
	err = stream.Send(&proto.RunRequest{
		Options: &proto.RunOptions{
			Parser:    "list",
			Formatter: "json",
		},
		Input: []byte("안녕 hello\n"),
	})
	if err != nil {
		panic(err)
	}

	t.Log("stream.Send[1]")
	err = stream.Send(&proto.RunRequest{
		Input: []byte("고양이 cat\n"),
	})
	if err != nil {
		panic(err)
	}

	t.Log("stream.CloseSend")
	err = stream.CloseSend()
	if err != nil {
		panic(err)
	}

	t.Log("waitc.block")
	<-waitc
}
