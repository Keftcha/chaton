package types

import (
	"github.com/keftcha/chaton/grpc/chaton"
)

// Client represent a client connetion to the server
type Client struct {
	Stream chaton.Chaton_JoinServer
	Nick   string
	Status string
}
