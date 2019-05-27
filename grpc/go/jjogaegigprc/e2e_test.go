package jjogaegigprc

import (
	"bytes"
	"context"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(1000)

	go func() {
		s := grpc.NewServer()
		RegisterRunServiceServer(s, NewServer())
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func newClient() (RunServiceClient, func()) {
	dialer := grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	})

	conn, err := grpc.DialContext(context.Background(), "", dialer, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return NewRunServiceClient(conn), func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}
}

func TestRun(t *testing.T) {
	c, teardown := newClient()
	defer teardown()
	assertRun(t, &SimpleRunner{c})
}

func TestRunStream(t *testing.T) {
	c, teardown := newClient()
	defer teardown()
	assertRun(t, &StreamingRunner{c})
}

func assertRun(t *testing.T, runner Runner) {
	inBuf := bytes.NewBufferString("안녕 hello\n고양이 cat\n")
	outBuf := &bytes.Buffer{}

	err := runner.Run(
		inBuf,
		outBuf,
		&RunConfig{
			Parser:    "list",
			Formatter: "csv",
			Options: map[string]string{
				pkg.OPT_HEADER: "HEADER",
				pkg.OPT_DEBUG:  "true",
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, outBuf.String(), "HEADER")
	assert.Contains(t, outBuf.String(), ",,안녕,,,hello,,,,,,,,,")
	assert.Contains(t, outBuf.String(), ",,고양이,,,cat,,,,,,,,,")
}
