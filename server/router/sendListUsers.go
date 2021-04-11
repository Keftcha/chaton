package router

import (
	"github.com/keftcha/chaton/grpc/chaton"
	"github.com/keftcha/chaton/server/types"
)

func SendListUsers(cs types.Clients, e types.Event) {
	msg := cs.ListClients()

	// Send only to the user who ask who is here
	e.Client.Stream.Send(
		&chaton.Event{
			Type: chaton.MsgType_MESSAGE,
			Msg: &chaton.Msg{
				Content: msg,
			},
		},
	)
}
