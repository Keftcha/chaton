package main

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

// Event represent a grpc event and the client that send it
type Event struct {
	Client *Client
	Event  *chaton.Event
}

// NewEventWithoutClient create an event without client
func NewEventWithoutClient(
	eventType chaton.MsgType,
	msgContent string,
) Event {
	return Event{
		Event: &chaton.Event{
			Type: eventType,
			Msg: &chaton.Msg{
				Content: msgContent,
			},
		},
		Client: &Client{
			Stream: nil,
			Nick:   "",
			Status: "",
		},
	}
}
