package main

//go:generate protoc --go_out=../grpc/ --go-grpc_out=../grpc/ -I ../proto/ ../proto/chaton.proto

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/keftcha/chaton/grpc/chaton"
)

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
	go routeEvents(c)
	return s
}

// Connect implements the chaton interface
func (s *ChatonServer) Join(stream chaton.Chaton_JoinServer) error {
	// Initialise the client of the event
	e := event{
		client: &client{
			stream: stream,
			nick:   "",
			status: "Online.",
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

func routeEvents(c <-chan event) {
	var clients clients

	// Loop on messages in channel
	for e := range c {
		// Switch on the event
		switch e.event.Type {
		// A new client has arrived
		case chaton.MsgType_CONNECT:
			connect(&clients, e)
		// Client change his nickname
		case chaton.MsgType_SET_NICKNAME:
			changeNick(clients, e)
		// Client send message
		case chaton.MsgType_MESSAGE:
			clients.broadcasting(e)
		// Client whant to leave us
		case chaton.MsgType_QUIT:
			quit(&clients, e)
		// Client do an action
		case chaton.MsgType_ME:
			action(clients, e)
		// List users on the server
		case chaton.MsgType_LIST:
			sendListUsers(clients, e)
		}
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
