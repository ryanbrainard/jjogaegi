package jjogaegigprc

import (
	"bytes"
	"context"
	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/run"
	"io"
	"log"
	"sync"
)

type server struct{}

func NewServer() RunServiceServer {
	return &server{}
}

func (s *server) Run(ctx context.Context, req *RunRequest) (*RunResponse, error) {
	output := &bytes.Buffer{}

	err := run.Run(
		bytes.NewReader(req.Input),
		output,
		cmd.ParseOptParser(req.Config.Parser),
		cmd.ParseOptFormatter(req.Config.Formatter),
		req.Config.Options,
	)

	return &RunResponse{Output: output.Bytes()}, err
}

func (s *server) RunStream(stream RunService_RunStreamServer) error {
	log.Println("fn=RunStream run_stream.start")

	inputBuf := newBlockingReadWriteCloser()
	runOnce := sync.Once{}
	waitc := make(chan struct{})

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("fn=RunStream stream.receive.eof")
				inputBuf.Close()
				break
			}

			log.Printf("fn=RunStream stream.receive.err %v", err)
			return err
		}

		log.Printf("fn=RunStream stream.receive.write [%v]", string(req.Input))
		if _, err := inputBuf.Write(req.Input); err != nil {
			log.Printf("fn=RunStream stream.receive.write.err %v", err)
			return err
		}

		runOnce.Do(func() {
			go func() {
				log.Println("fn=RunStream run.start")

				if err := run.Run(
					inputBuf,
					&streamWriter{stream},
					cmd.ParseOptParser(req.Config.Parser),
					cmd.ParseOptFormatter(req.Config.Formatter),
					req.Config.Options,
				); err != nil {
					panic(err)
				}

				log.Println("fn=RunStream run.done")
				close(waitc)
			}()
		})
	}

	log.Println("fn=RunStream run_stream.wait")
	<-waitc
	log.Println("fn=RunStream run_stream.done")
	return nil
}

type blockingReadWriteCloser struct {
	ch chan byte
}

func newBlockingReadWriteCloser() io.ReadWriteCloser {
	return &blockingReadWriteCloser{ch: make(chan byte, 1024)}
}

func (r *blockingReadWriteCloser) Write(p []byte) (int, error) {
	for _, b := range p {
		r.ch <- b
	}
	return len(p), nil
}

func (r *blockingReadWriteCloser) Read(p []byte) (n int, err error) {
	i := 0
	for ; i < len(p); i++ {
		select {
		case b := <-r.ch:
			// non-blocking: try to read as long as we have bytes in the channel
			p[i] = b
		default:
			// blocking: read at least one byte
			if i == 0 {
				b := <-r.ch
				p[i] = b
			}
		}
	}
	return i, nil
}

func (r *blockingReadWriteCloser) Close() (err error) {
	close(r.ch)
	return nil
}

type streamWriter struct {
	stream RunService_RunStreamServer
}

func (sw *streamWriter) Write(p []byte) (int, error) {
	log.Printf("fn=RunStream streamWriter.send [%v]", string(p))
	err := sw.stream.Send(&RunResponse{Output: p})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
