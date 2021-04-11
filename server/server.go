package main

//go:generate protoc --go_out=../grpc/ --go-grpc_out=../grpc/ -I ../proto/ ../proto/chaton.proto

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/keftcha/chaton/grpc/chaton"
	"github.com/keftcha/chaton/server/router"
	"github.com/keftcha/chaton/server/types"
)

// ChatonServer implements the Chaton service
type ChatonServer struct {
	chaton.UnimplementedChatonServer

	// Events sended
	es chan<- types.Event
}

// newChatonServer create a new ChatonServer service
func newChatonServer() *ChatonServer {
	c := make(chan types.Event)
	s := &ChatonServer{
		es: c,
	}
	go routeEvents(c)
	return s
}

// Join implements the chaton interface
func (s *ChatonServer) Join(stream chaton.Chaton_JoinServer) error {
	// Initialise the client of the event
	e := types.Event{
		Client: &types.Client{
			Stream: stream,
			Nick:   "",
			Status: "Online.",
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
		e.Event = in
		s.es <- e
	}
}

func routeEvents(c <-chan types.Event) {
	var clients types.Clients

	// Loop on messages in channel
	for e := range c {
		// Switch on the Event
		switch e.Event.Type {
		// A new client has arrived
		case chaton.MsgType_CONNECT:
			router.Connect(&clients, e)
		// Client change his nickname
		case chaton.MsgType_SET_NICKNAME:
			router.ChangeNick(clients, e)
		// Client send message
		case chaton.MsgType_MESSAGE:
			clients.Broadcasting(e)
		// Client whant to leave us
		case chaton.MsgType_QUIT:
			router.Quit(&clients, e)
		// Client do an action
		case chaton.MsgType_ME:
			router.Action(clients, e)
		// List users on the server
		case chaton.MsgType_LIST:
			router.SendListUsers(clients, e)
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
