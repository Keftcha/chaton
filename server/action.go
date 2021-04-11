package main

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
)

func action(cs clients, e event) {
	cs.broadcasting(
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
}
