package main

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
)

func changeNick(cs clients, e event) {
	// The new nickname is the content of the message
	newNick := e.event.Msg.Content
	// Prevent other users that the client has changed his nickname
	cs.broadcasting(
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
}
