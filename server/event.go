package main

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

// event represent a grpc event and the client that send it
type event struct {
	client *client
	event  *chaton.Event
}

// newEvent create an event
func newEvent(
	eventType chaton.MsgType,
	msgContent string,
	stream chaton.Chaton_JoinServer,
	nick string,
) event {
	return event{
		event: &chaton.Event{
			Type: eventType,
			Msg: &chaton.Msg{
				Content: msgContent,
			},
		},
		client: &client{
			stream: stream,
			nick:   nick,
		},
	}
}
