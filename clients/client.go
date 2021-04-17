package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/keftcha/chaton/grpc/chaton"
)

func join(c chaton.ChatonClient) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := c.Join(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	// Send connection event
	stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_CONNECT,
			Msg: &chaton.Msg{
				Content: "Moritz",
			},
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

			fmt.Printf(
				"%s: %s\n",
				recv.Msg.Author,
				recv.Msg.Content,
			)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		content, _ := reader.ReadString('\n')

		content = strings.TrimSpace(content)
		// Split the content message in words
		msg := strings.Split(content, " ")

		if len(msg) == 0 {
			continue
		}
		msgType := chaton.MsgType_MESSAGE
		switch msg[0] {
		case "/nick":
			msgType = chaton.MsgType_SET_NICKNAME
			content = strings.Join(msg[1:], " ")
		case "/me":
			msgType = chaton.MsgType_ME
			content = strings.Join(msg[1:], " ")
		case "/list":
			msgType = chaton.MsgType_LIST
		case "/quit":
			msgType = chaton.MsgType_QUIT
			content = strings.Join(msg[1:], " ")
		case "/status":
			msgType = chaton.MsgType_STATUS
			content = strings.Join(msg[1:], " ")
		case "/clear":
			msgType = chaton.MsgType_CLEAR
		case "/show":
			msgType = chaton.MsgType_SHOW
		}

		err := stream.Send(
			&chaton.Event{
				Type: msgType,
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

	join(client)
}
