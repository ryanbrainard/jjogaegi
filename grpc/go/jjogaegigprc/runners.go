package jjogaegigprc

import (
	"bufio"
	"context"
	"io"
	"io/ioutil"
	"log"
)

// Runner is an interface of wrappers around the generated client
type Runner interface {
	Run(in io.Reader, out io.Writer, config *RunConfig) error
}

type simpleRunner struct {
	client RunServiceClient
}

func (r *simpleRunner) Run(in io.Reader, out io.Writer, config *RunConfig) error {
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return nil
	}

	response, err := r.client.Run(context.TODO(), &RunRequest{
		Config: config,
		Input:  input,
	})
	if err != nil {
		return nil
	}

	_, err = out.Write(response.Output)
	return err
}

type streamingRunner struct {
	client RunServiceClient
}

func (r *streamingRunner) Run(in io.Reader, out io.Writer, config *RunConfig) error {
	stream, err := r.client.RunStream(context.TODO())
	if err != nil {
		return err
	}

	waitc := make(chan struct{})

	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("fn=streamingRunner.Run at=stream.receive.eof")
					close(waitc)
					return
				}
				panic(err) // TODO
			}

			log.Printf("fn=streamingRunner.Run at=stream.receive.output [%v]", string(response.Output))
			if _, err := out.Write(response.Output); err != nil {
				panic(err) // TODO
			}
		}
	}()

	inScan := bufio.NewScanner(in)
	inScan.Split(bufio.ScanLines)
	for inScan.Scan() {
		err := stream.Send(&RunRequest{
			Config: config,                                  // only actually used server-side on the first req
			Input:  append(inScan.Bytes(), []byte("\n")...), // re-append after consumed by bufio.ScanLines
		})
		if err != nil {
			return err
		}
	}

	log.Println("fn=streamingRunner.Run at=stream.CloseSend")
	if err := stream.CloseSend(); err != nil {
		return err
	}

	<-waitc
	return nil
}
