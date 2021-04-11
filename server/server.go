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
	stream chaton.Chaton_JoinServer
	nick   string
	status string
}

type clients []*client

func (cs *clients) broadcasting(e event) {
	for _, c := range *cs {
		// Put the sender name as the message author
		e.event.Msg.Author = e.client.nick
		// Remove client if there is an error sending him message
		if err := c.stream.Send(e.event); err != nil {
			cs.remove(e.client)
		}
	}
}

func (cs *clients) remove(client *client) {
	// Remove the client of our list
	for i, c := range *cs {
		if c == client {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			break
		}
	}
}

func (cs *clients) listClients() string {
	lst := ""
	for i, c := range *cs {
		if i != 0 {
			lst += "\n"
		}
		lst += fmt.Sprintf(
			"- %s: %s",
			c.nick,
			c.status,
		)
	}
	return lst
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

func sendListUsers(cs clients, e event) {
	msg := cs.listClients()

	// Send only to the user who ask who is here
	e.client.stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_MESSAGE,
			Msg: &chaton.Msg{
				Content: msg,
			},
		},
	)
}

func routeEvents(c <-chan event) {
	var clients clients

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
			clients.broadcasting(
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
			clients.broadcasting(
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
			clients.broadcasting(e)
		// Client whant to leave us
		case chaton.MsgType_QUIT:
			// Did the client let a reason
			reason := ""
			if e.event.Msg != nil {
				reason = fmt.Sprintf(" (\"%s\")", e.event.Msg.Content)
			}
			// Prevent other users the client has left
			clients.broadcasting(
				event{
					event: &chaton.Event{
						Type: chaton.MsgType_MESSAGE,
						Msg: &chaton.Msg{
							Content: fmt.Sprintf(
								"%s has left",
								e.client.nick,
							) + reason,
						},
					},
					client: &client{
						stream: nil,
						nick:   "Server say",
					},
				},
			)
			// Remove the client of our list
			clients.remove(e.client)
		// Client do an action
		case chaton.MsgType_ME:
			clients.broadcasting(
				event{
					event: &chaton.Event{
						Type: chaton.MsgType_MESSAGE,
						Msg: &chaton.Msg{
							// Add the pseudo before the action
							Content: fmt.Sprintf(
								"*%s %s*",
								e.client.nick,
								e.event.Msg.Content,
							),
						},
					},
					client: &client{
						stream: nil,
						nick:   "Server say",
					},
				},
			)
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
