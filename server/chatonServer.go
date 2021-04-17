package main

import (
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/keftcha/chaton/grpc/chaton"
)

// ChatonServer implements the Chaton service
type ChatonServer struct {
	chaton.UnimplementedChatonServer

	// Events sended
	es chan event
	// Clients connected
	cs clients
}

// newChatonServer create a new ChatonServer service
func newChatonServer() *ChatonServer {
	s := &ChatonServer{
		es: make(chan event),
	}
	go s.dispatch()
	return s
}

// Join implements the chaton interface
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

// dispatch recieved event
func (s *ChatonServer) dispatch() {

	// Loop on messages in channel
	for e := range s.es {
		// Switch on the Event
		switch e.event.Type {
		// A new client has arrived
		case chaton.MsgType_CONNECT:
			s.connect(e)
		// Client change his nickname
		case chaton.MsgType_SET_NICKNAME:
			s.changeNick(e)
		// Client send message
		case chaton.MsgType_MESSAGE:
			s.cs.broadcasting(e)
		// Client whant to leave us
		case chaton.MsgType_QUIT:
			s.quit(e)
		// Client do an action
		case chaton.MsgType_ME:
			s.action(e)
		// List users on the server
		case chaton.MsgType_LIST:
			s.listUsers(e)
		// Client set his status
		case chaton.MsgType_STATUS:
			s.changeStatus(e)
		// Client remove his status
		case chaton.MsgType_CLEAR:
			s.clearStatus(e)
		// Show to the client his status
		case chaton.MsgType_SHOW:
			s.showStatus(e)
		}
	}
}

// connect a client to the server and send him recieved events
func (s *ChatonServer) connect(e event) {
	// Set the client name to the content of the message
	if e.event.Msg != nil {
		e.client.nick = e.event.Msg.Content
	} else {
		id, _ := uuid.NewRandom()
		e.client.nick = id.String()
	}

	// Add the client to our list
	s.cs.add(e.client)

	// Prevent other users that a new client has arrived
	s.cs.broadcasting(
		newEvent(
			chaton.MsgType_CONNECT,
			fmt.Sprintf("%s has joined.", e.client.nick),
		),
	)
}

// changeNick name of a client
func (s *ChatonServer) changeNick(e event) {
	// The new nickname is the content of the message
	newNick := e.event.Msg.Content

	// Prevent other users that the client has changed his nickname
	s.cs.broadcasting(
		newEvent(
			chaton.MsgType_SET_NICKNAME,
			fmt.Sprintf("%s is now known as %s", e.client.nick, newNick),
		),
	)

	// Change the nickname
	e.client.nick = newNick
}

// quit remove a connected client
func (s *ChatonServer) quit(e event) {
	// Did the client let a reason
	reason := ""
	if e.event.Msg != nil {
		reason = fmt.Sprintf(" (\"%s\")", e.event.Msg.Content)
	}

	// Prevent other users the client has left
	s.cs.broadcasting(
		newEvent(
			chaton.MsgType_QUIT,
			fmt.Sprintf("%s has left.", e.client.nick)+reason,
		),
	)

	// Remove the client of our list
	s.cs.remove(e.client)
}

// action a client made
func (s *ChatonServer) action(e event) {
	s.cs.broadcasting(
		newEvent(
			chaton.MsgType_ME,
			// Add the pseudo before the action
			fmt.Sprintf("%s %s", e.client.nick, e.event.Msg.Content),
		),
	)
}

// listUsers connected to the server
func (s *ChatonServer) listUsers(e event) {
	msg := s.cs.listClients()

	// Send only to the user who ask who is here
	e.client.stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_LIST,
			Msg: &chaton.Msg{
				Content: msg,
			},
		},
	)
}

// changeStatus let the user set his status
func (s *ChatonServer) changeStatus(e event) {
	e.client.status = e.event.Msg.Content
}

// clearStatus of the user
func (s *ChatonServer) clearStatus(e event) {
	e.client.status = "Online."
}

// showStatus send the status to the client
func (s *ChatonServer) showStatus(e event) {
	e.client.stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_SHOW,
			Msg: &chaton.Msg{
				Content: "Current status: " + e.client.status,
			},
		},
	)
}
