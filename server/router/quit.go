package router

import (
	"fmt"

	"github.com/keftcha/chaton/grpc/chaton"
	"github.com/keftcha/chaton/server/types"
)

func Quit(cs *types.Clients, e types.Event) {
	// Did the client let a reason
	reason := ""
	if e.Event.Msg != nil {
		reason = fmt.Sprintf(" (\"%s\")", e.Event.Msg.Content)
	}
	// Prevent other users the client has left
	cs.Broadcasting(
		types.Event{
			Event: &chaton.Event{
				Type: chaton.MsgType_MESSAGE,
				Msg: &chaton.Msg{
					Content: fmt.Sprintf(
						"%s has left",
						e.Client.Nick,
					) + reason,
				},
			},
			Client: &types.Client{
				Stream: nil,
				Nick:   "Server say",
			},
		},
	)
	// Remove the client of our list
	cs.Remove(e.Client)
}
