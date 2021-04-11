package main

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

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
