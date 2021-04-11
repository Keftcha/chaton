package types

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

// Event represent a grpc event and the client that send it
type Event struct {
	Client *Client
	Event  *chaton.Event
}
