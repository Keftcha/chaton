package main

//go:generate protoc --go_out=../grpc/ --go-grpc_out=../grpc/ -I ../proto/ ../proto/chaton.proto

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	"github.com/keftcha/chaton/grpc/chaton"
)

type event struct {
	client *client
	event  *chaton.Event
}

type client struct {
	stream chaton.Chaton_ConnectServer
	nick   string
}

// ChatonServer implements the Chaton service
type ChatonServer struct {
	chaton.UnimplementedChatonServer

	// Events sended
	es chan<- event
}

// newChatonServer create a new ChatonServer service
func newChatonServer() *ChatonServer {
	c := make(chan event)
	s := &ChatonServer{
		es: c,
	}
	go eventRouting(c)
	return s
}

// Connect implements the chaton interface
func (s *ChatonServer) Connect(stream chaton.Chaton_ConnectServer) error {
	// Initialise the client of the event
	e := event{
		client: &client{
			stream: stream,
			nick:   "",
		},
	}

	// Infinite loop that recieve messages
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Add the new message to the channel queue
		e.event = in
		s.es <- e
	}
}

func eventRouting(c <-chan event) {
	var clients []*client

	// Loop on messages in channel
	for e := range c {
		// Switch on the event
		switch e.event.Type {
		// A new client has arrived
		case chaton.MsgType_CONNECT:
			// Set the client name to the content of the message
			if e.event.Msg != nil {
				e.client.nick = e.event.Msg.Content
			} else {
				id, _ := uuid.NewRandom()
				e.client.nick = id.String()
			}
			clients = append(clients, e.client)
			// Prevent other users that a new client has arrived
			broadcasting(
				clients,
				event{
					event: &chaton.Event{
						Type: chaton.MsgType_MESSAGE,
						Msg: &chaton.Msg{
							Content: "A new boi has arrived",
						},
					},
					client: &client{
						stream: nil,
						nick:   "Server say",
					},
				},
			)
		// Client change his nickname
		case chaton.MsgType_SET_NICKNAME:
			// The new nickname is the content of the message
			newNick := e.event.Msg.Content
			// Prevent other users that the client has changed his nickname
			broadcasting(
				clients,
				event{
					event: &chaton.Event{
						Type: chaton.MsgType_MESSAGE,
						Msg: &chaton.Msg{
							Content: fmt.Sprintf(
								"%s change his nickname to %s",
								e.client.nick,
								newNick,
							),
						},
					},
					client: &client{
						stream: nil,
						nick:   "Server say",
					},
				},
			)
			// Change the nickname
			e.client.nick = newNick
		// Client send message
		case chaton.MsgType_MESSAGE:
			broadcasting(clients, e)
		case chaton.MsgType_QUIT:
		case chaton.MsgType_ME:
		case chaton.MsgType_LIST:
		}

	}
}

func broadcasting(cs []*client, e event) {
	for _, c := range cs {
		// Put the sender name as the message author
		e.event.Msg.Author = e.client.nick
		c.stream.Send(e.event)
	}
}

func main() {
	// Initialise the listening connections
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 21617))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server started on port " + lis.Addr().String())

	// Create and start the gRPC server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chaton.RegisterChatonServer(grpcServer, newChatonServer())
	grpcServer.Serve(lis)
}
