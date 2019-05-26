package server

import (
	"bytes"
	"context"
	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/grpc/proto"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/ryanbrainard/jjogaegi/run"
	"io"
	"log"
	"sync"
)

type server struct {
}

func NewServer() proto.JjogaegiServer {
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

type streamWriter struct {
	stream proto.Jjogaegi_RunStreamServer
}

func (sw *streamWriter) Write(p []byte) (int, error) {
	//log.Printf("streamWriter.Write %v", string(p))

	err := sw.stream.Send(&proto.RunResponse{Output: p})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
