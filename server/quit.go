package main

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
)

func quit(cs *clients, e event) {
	// Did the client let a reason
	reason := ""
	if e.event.Msg != nil {
		reason = fmt.Sprintf(" (\"%s\")", e.event.Msg.Content)
	}
	// Prevent other users the client has left
	cs.broadcasting(
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
	cs.remove(e.client)
}
