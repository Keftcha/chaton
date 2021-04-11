package router

import (
	"github.com/google/uuid"

	"github.com/keftcha/chaton/grpc/chaton"
	"github.com/keftcha/chaton/server/types"
)

func Connect(cs *types.Clients, e types.Event) {
	// Set the client name to the content of the message
	if e.Event.Msg != nil {
		e.Client.Nick = e.Event.Msg.Content
	} else {
		id, _ := uuid.NewRandom()
		e.Client.Nick = id.String()
	}
	// Add the client to our list
	cs.Add(e.Client)
	// Prevent other users that a new client has arrived
	cs.Broadcasting(
		types.Event{
			Event: &chaton.Event{
				Type: chaton.MsgType_MESSAGE,
				Msg: &chaton.Msg{
					Content: "A new boi has arrived",
				},
			},
			Client: &types.Client{
				Stream: nil,
				Nick:   "Server say",
			},
		},
	)
}
