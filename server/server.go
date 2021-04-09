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

type event struct {
	client chaton.Chaton_ConnectServer
	event  *chaton.Event
}

// ChatonServer implements the Chaton service
type ChatonServer struct {
	chaton.UnimplementedChatonServer

	// Events sended
	es chan<- event
	// Client streams
	streams []chaton.Chaton_ConnectServer
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
		client: stream,
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
	clients := make([]chaton.Chaton_ConnectServer, 0)

	// Loop on messages in channel
	for e := range c {
		// Switch on the event
		switch e.event.Type {
		// A new client has arrived
		case chaton.MsgType_CONNECT:
			clients = append(clients, e.client)
			broadcasting(
				clients,
				&chaton.Event{
					Type: chaton.MsgType_MESSAGE,
					Msg: &chaton.Msg{
						Content: "A new boi has arrived",
					},
				},
			)
			continue
		case chaton.MsgType_SET_NICKNAME:
		case chaton.MsgType_MESSAGE:
		case chaton.MsgType_QUIT:
		case chaton.MsgType_ME:
		case chaton.MsgType_LIST:
		}

	}
}

func broadcasting(cs []chaton.Chaton_ConnectServer, e *chaton.Event) {
	for _, c := range cs {
		c.Send(e)
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
