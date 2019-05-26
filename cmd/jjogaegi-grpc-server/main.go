package main

import (
	"bytes"
	"context"
	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/ryanbrainard/jjogaegi/proto"
	"github.com/ryanbrainard/jjogaegi/run"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
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

func (s *server) Run(ctx context.Context, req *proto.RunRequest) (*proto.RunResponse, error) {
	output := &bytes.Buffer{}

	err := run.Run(
		bytes.NewReader(req.Input),
		output,
		cmd.ParseOptParser(req.Options.Parser),
		cmd.ParseOptFormatter(req.Options.Formatter),
		map[string]string{},
	)

	return &proto.RunResponse{Output: output.Bytes()}, err
}

func (s *server) RunStream(stream proto.Jjogaegi_RunStreamServer) error {
	log.Println("run_stream.start")

	inputBuf := &bytes.Buffer{}
	runOne := sync.Once{}
	waitc := make(chan struct{})

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("stream.receive.eof")
				break
			}

			log.Printf("stream.receive.err %v", err)
			return err
		}

		log.Printf("stream.receive.write.start")
		if _, err := inputBuf.Write(req.Input); err != nil {
			log.Printf("stream.receive.write.err %v", err)
			return err
		}
		log.Printf("stream.receive.write.done")

		runOne.Do(func() {
			go func() {
				log.Println("run.start")
				if err := run.Run(
					inputBuf,
					&streamWriter{stream},
					cmd.ParseOptParser(req.Options.Parser),
					cmd.ParseOptFormatter(req.Options.Formatter),
					map[string]string{
						pkg.OPT_DEBUG: "true",
					},
				); err != nil {
					panic(err)
				}
				log.Println("run.done")
				close(waitc)
			}()
		})
	}

	log.Println("run_stream.wait")
	<-waitc
	log.Println("run_stream.done")
	return nil
}

//type streamReader struct {
//	stream proto.Jjogaegi_RunStreamServer
//	buf *bytes.Buffer
//}
//
//func (sr *streamReader) Read(p []byte) (int, error) {
//	log.Println("streamReader.Read")
//
//	req, err := sr.stream.Recv()
//	if err != nil {
//		return 0, err
//	}
//	sr.buf.Write(req.Input)
//
//	return sr.buf.Read(p)
//}

type streamWriter struct {
	stream proto.Jjogaegi_RunStreamServer
}

func (sw *streamWriter) Write(p []byte) (int, error) {
	log.Printf("streamWriter.Write %v", string(p))

	err := sw.stream.Send(&proto.RunResponse{Output: p})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
