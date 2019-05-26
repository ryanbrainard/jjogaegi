package main

import (
	"context"
	"fmt"
	"github.com/ryanbrainard/jjogaegi/proto"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewJjogaegiClient(conn)

	// RUN TEST

	//response, err := c.Run(context.TODO(), &proto.RunRequest{
	//	Options: &proto.RunOptions{
	//		Parser:    "list",
	//		Formatter: "json",
	//	},
	//	Input:     []byte("안녕 hello\n고양이 cat"),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("output \n%+v", string(response.Output))

	// RUN STREAM TEST

	stream, err := c.RunStream(context.TODO())
	if err != nil {
		panic(err)
	}
	waitc := make(chan struct{})

	go func() {
		for {
			log.Println("waiting for response")
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("waitc.close")
					close(waitc)
					return
				}
				panic(fmt.Sprintf("-> %+v", err))
			}
			log.Printf("output (%d) \n%+v", len(response.Output), string(response.Output))
		}
	}()

	log.Println("stream.Send[0]")
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

	log.Println("stream.Send[1]")
	err = stream.Send(&proto.RunRequest{
		Input: []byte("고양이 cat\n"),
	})
	if err != nil {
		panic(err)
	}

	log.Println("stream.CloseSend")
	err = stream.CloseSend()
	if err != nil {
		panic(err)
	}

	log.Println("waitc.block")
	<-waitc
}
