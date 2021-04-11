package router

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
	"github.com/keftcha/chaton/server/types"
)

func ChangeNick(cs types.Clients, e types.Event) {
	// The new nickname is the content of the message
	newNick := e.Event.Msg.Content
	// Prevent other users that the client has changed his nickname
	cs.Broadcasting(
		types.Event{
			Event: &chaton.Event{
				Type: chaton.MsgType_MESSAGE,
				Msg: &chaton.Msg{
					Content: fmt.Sprintf(
						"%s change his nickname to %s",
						e.Client.Nick,
						newNick,
					),
				},
			},
			Client: &types.Client{
				Stream: nil,
				Nick:   "Server say",
			},
		},
	)
	// Change the nickname
	e.Client.Nick = newNick
}
