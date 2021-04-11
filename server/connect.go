package main

import (
	"github.com/google/uuid"

	"github.com/keftcha/chaton/grpc/chaton"
)

func connect(cs *clients, e event) {
	// Set the client name to the content of the message
	if e.event.Msg != nil {
		e.client.nick = e.event.Msg.Content
	} else {
		id, _ := uuid.NewRandom()
		e.client.nick = id.String()
	}
	// Add the client to our list
	cs.add(e.client)
	// Prevent other users that a new client has arrived
	cs.broadcasting(
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
}
