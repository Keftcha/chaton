package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

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
				Content: uname,
			},
		},
	)

	go func() {
		for {
			recv, err := stream.Recv()
			// Ignore EOF error
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}

			// Format time
			now := time.Now()
			displayedTime := fmt.Sprintf(
				"%02d\033[33m:\033[0m%02d\033[33m:\033[0m%02d",
				now.Hour(),
				now.Minute(),
				now.Second(),
			)

			// Format sender and message content
			author := recv.Msg.Author
			content := recv.Msg.Content
			switch recv.Type {
			case chaton.MsgType_CONNECT:
				author = "       \033[32m-->\033[0m"
			case chaton.MsgType_SET_NICKNAME:
				author = "        \033[35m--\033[0m"
			case chaton.MsgType_MESSAGE:
				author = "@" + recv.Msg.Author
			case chaton.MsgType_QUIT:
				author = "       \033[31m<--\033[0m"
			case chaton.MsgType_ME:
				author = "*"
				content = "\033[03m" + content + "\033[0m"
			case chaton.MsgType_LIST:
				author = "        \033[35m--\033[0m"
				// New message content
				c := ""
				for i, l := range strings.Split(content, "\n") {
					if i == 0 {
						c += l
					} else {
						c += "\n                    \033[32m|\033[0m  " + l
					}
				}
				content = c
			case chaton.MsgType_SHOW:
				author = "        \033[35m--\033[0m"
			}

			fmt.Printf(
				"%s %10s \033[32m|\033[0m  %s\n",
				displayedTime,
				author,
				content,
			)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
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
			log.Fatal(err)
		}

		// Properly close client stream
		if msg[0] == "/quit" {
			time.Sleep(100 * time.Millisecond)
			stream.CloseSend()
			return
		}
	}
}

var host string = "localhost"
var port int = 21617
var uname string // User nick name

func init() {
	flag.StringVar(&host, "host", "localhost", "The host address of the server")
	flag.IntVar(&port, "port", 21617, "The port of the server we connect to")
	flag.StringVar(&uname, "username", "", "The nickname of the user")

	flag.Parse()
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := chaton.NewChatonClient(conn)

	join(client)
}
