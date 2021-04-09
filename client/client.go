package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/keftcha/chaton/grpc/chaton"
)

func connect(c chaton.ChatonClient) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := c.Connect(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	// Send connection event
	stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_CONNECT,
			Msg:  nil,
		},
	)

	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Connection closed by server (EOF)")
				return
			}
			if err != nil {
				fmt.Println("err recieved message: ", err)
				return
			}

			fmt.Println(recv.Msg.Author, recv.Msg.Content)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		content, _ := reader.ReadString('\n')

		err := stream.Send(
			&chaton.Event{
				Type: chaton.MsgType_MESSAGE,
				Msg:  &chaton.Msg{Content: content},
			},
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:21617", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := chaton.NewChatonClient(conn)

	connect(client)
}
