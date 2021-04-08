package main

//go:generate protoc --go_out=../grpc/ --go-grpc_out=../grpc/ -I ../proto ../proto/chaton.proto

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	chaton "github.com/keftcha/chaton/grpc"
)

type message struct {
	client chaton.Chaton_ConnectServer
	msg    *chaton.Msg
}

// ChatonServer implements the Chaton service
type ChatonServer struct {
	chaton.UnimplementedChatonServer

	// Messages sended
	msgs chan<- message
	// Client streams
	streams []chaton.Chaton_ConnectServer
}

// newChatonServer create a new ChatonServer service
func newChatonServer() *ChatonServer {
	c := make(chan message)
	s := &ChatonServer{
		msgs: c,
	}
	go broadcasting(c)
	return s
}

// Connect implements the chaton interface
func (s *ChatonServer) Connect(stream chaton.Chaton_ConnectServer) error {
	m := message{
		client: stream,
		msg:    &chaton.Msg{Content: "¤$£"},
	}
	s.msgs <- m

	// Infinite loop that recieve messages
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Add the new message to the channel queue
		m.msg = msg
		s.msgs <- m
	}
}

func broadcasting(c <-chan message) {
	clients := make([]chaton.Chaton_ConnectServer, 0)
	_ = clients

	// Loop on messages in channel
	for m := range c {
		// A new client has arrived
		if m.msg.Content == "¤$£" {
			clients = append(clients, m.client)
			continue
		}
		for _, c := range clients {
			c.Send(m.msg)
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
