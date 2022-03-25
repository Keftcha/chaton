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
	es chan Event
	// Clients connected
	cs Clients
}

// newChatonServer create a new ChatonServer service
func newChatonServer() *ChatonServer {
	s := &ChatonServer{
		es: make(chan Event),
		cs: make(Clients),
	}
	go s.dispatch()
	return s
}

// Join implements the chaton interface
func (s *ChatonServer) Join(stream chaton.Chaton_JoinServer) error {
	// Initialise the client of the event
	e := Event{
		Client: &Client{
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

// dispatch recieved event
func (s *ChatonServer) dispatch() {

	// Loop on messages in channel
	for e := range s.es {
		// Switch on the Event
		switch e.Event.Type {
		// A new client has arrived
		case chaton.MsgType_CONNECT:
			s.connect(e)
		// Client change his nickname
		case chaton.MsgType_SET_NICKNAME:
			s.changeNick(e)
		// Client send message
		case chaton.MsgType_MESSAGE:
			Broadcasting(s.cs, e)
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
func (s *ChatonServer) connect(e Event) {
	// Set the client name to the content of the message
	if e.Event.Msg != nil {
		e.Client.Nick = e.Event.Msg.Content
	} else {
		id, _ := uuid.NewRandom()
		e.Client.Nick = id.String()
	}

	// Add the client to our list
	AddClient(s.cs, e.Client)

	// Prevent other users that a new client has arrived
	Broadcasting(
		s.cs,
		NewEventWithoutClient(
			chaton.MsgType_CONNECT,
			fmt.Sprintf("%s has joined.", e.Client.Nick),
		),
	)
}

// changeNick name of a client
func (s *ChatonServer) changeNick(e Event) {
	// The new nickname is the content of the message
	newNick := e.Event.Msg.Content

	// Prevent other users that the client has changed his nickname
	Broadcasting(
		s.cs,
		NewEventWithoutClient(
			chaton.MsgType_SET_NICKNAME,
			fmt.Sprintf("%s is now known as %s", e.Client.Nick, newNick),
		),
	)

	// Change the nickname
	e.Client.Nick = newNick
}

// quit remove a connected client
func (s *ChatonServer) quit(e Event) {
	// Did the client let a reason
	reason := ""
	if e.Event.Msg != nil {
		reason = fmt.Sprintf(" (\"%s\")", e.Event.Msg.Content)
	}

	// Prevent other users the client has left
	Broadcasting(
		s.cs,
		NewEventWithoutClient(
			chaton.MsgType_QUIT,
			fmt.Sprintf("%s has left.", e.Client.Nick)+reason,
		),
	)

	// Remove the client of our list
	RemoveClient(s.cs, e.Client)
}

// action a client made
func (s *ChatonServer) action(e Event) {
	Broadcasting(
		s.cs,
		NewEventWithoutClient(
			chaton.MsgType_ME,
			// Add the pseudo before the action
			fmt.Sprintf("%s %s", e.Client.Nick, e.Event.Msg.Content),
		),
	)
}

// listUsers connected to the server
func (s *ChatonServer) listUsers(e Event) {
	msg := ListClients(s.cs)

	// Send only to the user who ask who is here
	e.Client.Stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_LIST,
			Msg: &chaton.Msg{
				Content: msg,
			},
		},
	)
}

// changeStatus let the user set his status
func (s *ChatonServer) changeStatus(e Event) {
	e.Client.Status = e.Event.Msg.Content
}

// clearStatus of the user
func (s *ChatonServer) clearStatus(e Event) {
	e.Client.Status = "Online."
}

// showStatus send the status to the client
func (s *ChatonServer) showStatus(e Event) {
	e.Client.Stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_SHOW,
			Msg: &chaton.Msg{
				Content: "Current status: " + e.Client.Status,
			},
		},
	)
}
