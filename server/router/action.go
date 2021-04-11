package router

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
	"github.com/keftcha/chaton/server/types"
)

func Action(cs types.Clients, e types.Event) {
	cs.Broadcasting(
		types.Event{
			Event: &chaton.Event{
				Type: chaton.MsgType_MESSAGE,
				Msg: &chaton.Msg{
					// Add the pseudo before the action
					Content: fmt.Sprintf(
						"*%s %s*",
						e.Client.Nick,
						e.Event.Msg.Content,
					),
				},
			},
			Client: &types.Client{
				Stream: nil,
				Nick:   "Server say",
			},
		},
	)
}
